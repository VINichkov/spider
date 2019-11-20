package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"spider/block_list"
	"spider/models/entity"
	"spider/models/entity/client"
	"spider/models/entity/company"
	"spider/models/entity/job"
	"spider/proxy"
	"spider/source"
	"strings"
)

func getMainPage(location entity.Location, job_for_prepare chan entity.JobForPrepare,
	block *block_list.BlockList, source source.Source, thread int) {
	var (
		end          bool = false
		numberOfPage int
		iter         int
	)
Exit:
	for i := 0; i < MaxPage; i++ {
		if end || block.Include(location.Id) {
			break Exit
		}
		//Создаем УРЛ и запрашиваем страницу
		page := i + 1
		url := source.CreateQuery(location, page)


		resp, err := proxy.Connect(url, true)
		if err != nil {
			logging("getMainPage", thread,"Error " + err.Error())
			break Exit
		}

		//Преобразуем страницу в Nodes
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logging("getMainPage", thread,"Error " + err.Error())
			break Exit
		}

		resp.Body.Close()


		//Получаем общее количество страниц с 1 страницы
		if i==0 {
			numberOfAds, err := source.NumberOfAds(*doc)
			if err != nil {
				break Exit
			}
			//log.Info().Msg(strconv.Itoa(numberOfAds))
			if numberOfAds == 0 {
				break Exit
			}

			numberOfPage= numberOfAds / 10

			if numberOfAds % 10 > 0{
				numberOfPage += 1
			}

			if numberOfPage > MaxPage {
				iter = MaxPage
			} else {
				iter = numberOfPage
			}

		}

		nodes := doc.Find("div.result, a.result, div.job")
		/*log.Info().Msg("_____________________________________________________________________")
		 _ =nodes
		*/
		job_for_prepare <- entity.JobForPrepare{
			Job:      *nodes,
			Location: location,
			Page:     page,
		}
		nodes = nil
		if page == iter{
			break Exit
		}

	}
}


func getList(prepareJobs entity.JobForPrepare, gchGetJobs chan entity.GotJob,
	block *block_list.BlockList, source source.Source, thread int, rep *Crawler){
	//log.Debug().Msg(strconv.Itoa(prepareJobs.Job.Size()))

	//Если пришло 0 работ, то блокирует локацию и выходим
	if prepareJobs.Job.Size() == 0 {
		block.Push(prepareJobs.Location.Id)
		return
	}
	//Для всех работ
	prepareJobs.Job.Each(func(i int, job *goquery.Selection) {
		//Если локация заблокирована, то выходим
		if block.Include(prepareJobs.Location.Id){
			return
		}

		title := source.GetTitle(*job)
		//Если Наименование пусто, то выходим
		if title.Name == ""{
			log.Debug().Msg("Не определили наименование ваканcии")
			return
		}

		//Если старое, то выходим
		if how_long(source.AgeOfJob(*job)) == false {
			log.Debug().Msg("Старая вакансия")
			return
		}


		//Если не удалось получить URL, то выходим
		url, err := source.UrlToJob(title.Href)
		if err != nil {
			log.Debug().Msg(err.Error())
			return
		}


		//Если наименование компании пусто, то выходим
		company := source.Company(*job)
		if company==""{
			log.Debug().Msg("Пустая компания ")
			return
		}


		params := entity.NewCompareWithIndexIn(
			company,
			title.Name,
			url,
			prepareJobs.Location,
		)
		compareResult := rep.compareWithIndex(params)

		if compareResult.NoError(){

			//определяем з.п
			salaryMin := 0
			salaryMax := 0
			s := source.Salary(*job)
			if len(s) >= 2 {
				salaryMin = s[0]
				salaryMax = s[1]
			}
			if len(s) == 1 {
				salaryMin = s[0]
			}

			//Отправляем полученные данные в следующуу очередь

			gchGetJobs<- entity.GotJob{
				url,
				title.Name,
				company,
				compareResult.CompanyId,
				salaryMin,
				salaryMax,
				prepareJobs.Location,
				prepareJobs.Page,
			}
			return
		}

		//Найден подобный источник
		if compareResult.Is_SameSource(){
			log.Debug().Msg("Присуствует ссылка в БД. Закрываем локацию.")
			block.Push(prepareJobs.Id)
			return
		}

		//Найдо наименование и компания
		if compareResult.Is_SameTitle(){
			log.Debug().Msg("Присуствует наименование и компания. ")
			return
		}

		//В блок листе
		if compareResult.Is_BlockList(){
			log.Debug().Msg("Компания в блок листе")
			return
		}

		if compareResult.Is_AnyError(){
			log.Debug().Msg("Неопознанная ошибка")
			return
		}

	})
}


func get_job(got_job entity.GotJob, chJobsForSave chan entity.JobForSave,
	source source.Source,  thread int) {

	//Получаем страницу работы. Если ошибка, то выходим
	job, err := proxy.Connect(got_job.Link, true)
	if err != nil {
		logging("get_job", thread,"Error " + err.Error())
		return
	}

	//Преобразуем страницу в Nodes
	doc, err := goquery.NewDocumentFromReader(job.Body)
	if err != nil {
		logging("get_job", thread,"Error " + err.Error())
		return
	}


	defer job.Body.Close()

	//Получаем описание
	desc:= source.Description(*doc)
	deleteAttr(&desc)
	description, err := desc.Html()

	if err != nil{
		log.Debug().Msg("Описание работ нет")
		return
	}

	if description == ""{
		log.Debug().Msg("Нет описания вакансии")
		return
	}
	description = replace(&description)
	//Получаем apply ссылку
	apply := source.ApplyLinc(*doc, got_job)

	chJobsForSave <- entity.JobForSave{
		got_job,
		description,
		apply,
	}

}

func save_job(jobForSave entity.JobForSave, thread int, rep *Crawler){

	var flagOldCompany bool = true
	// Создаем компанию, потому что не пришла
	if jobForSave.CompanyId.Id == 0 {
		//Получили размер
		size, err := rep.sizes.First()
		if err != nil {
			logging("save_job", thread, err.Error() )
			return
		}

		newCompany := company.NewCompany(
			jobForSave.Company,
			size,
			jobForSave.Location.Id,
			"",
			"",
			false,
			"",
			0,
			false,
			"",
		)
		idComp, err := rep.company.Create(newCompany)
		if err != nil {
			logging("save_job", thread, err.Error() )
			return
		}

		jobForSave.CompanyId.Id = idComp
		flagOldCompany = false
	}

	var clientId  int = 0

	if flagOldCompany {
		clientId, _ = rep.clients.FindByCompanyIdFirst(jobForSave.CompanyId.Id)
	}

	if (!flagOldCompany || clientId == 0){
		newClient := client.NewSimulationClientForClawler(
			jobForSave.Company,
			jobForSave.CompanyId.Id,
			jobForSave.Location.Id,
		)

		createdClientId, err := rep.clients.Create(newClient)

		if err != nil {
			logging("save_job", thread, err.Error() )
			return
		}
		clientId = createdClientId
	}

	newJob := job.NewJob(
		strings.Title(jobForSave.Title),
		jobForSave.Location.Id,
		jobForSave.Salary_min,
		jobForSave.Salary_max,
		jobForSave.Description,
		jobForSave.CompanyId.Id,
		clientId,
		"",
		jobForSave.Link,
		jobForSave.Apply,
	)
	_, err := rep.job.Create(newJob)
	if err != nil {
		logging("save_job", thread, err.Error() )
		return
	}

}

