package crawler

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"spider/block_list"
	"spider/models"
	"spider/models/entity"
	"spider/models/repository"
	"spider/source"
	"strings"
	"sync"
)

const MaxPage int = 10

const GetListC = 4
const PrepareJobC = 1
const GetJobC = 14
const SaveJobC = 1

var wm sync.WaitGroup


type Crawler struct {
	locations models.LocationRepository
	company	models.CompanyRepository
	job models.JobRepository
	sizes models.SizeRepository
	clients models.ClientRepository

}


func NewCrawlerHandler(db *sqlx.DB) *Crawler  {
	return &Crawler{
		locations: repository.NewSQLLocationRepo(db),
		company:   repository.NewSQLCompanyRepo(db),
		job:       repository.NewSQLdbJobRepo(db),
		sizes:     repository.NewSQLdbSizeRepo(db),
		clients:   repository.NewSQLdbClientRepo(db),
	}
}

func (rep *Crawler)Start()  {
	log.Info().Msg("Start")
	wm.Add(2)
	go rep.StartForSource(source.InitJora())
	go rep.StartForSource(source.InitIndeed())
	wm.Wait()
	log.Info().Msg("End")
}

func (rep *Crawler)StartForSource(source source.Source){
	log.Info().Msg("Start Go")
	//Каналы для гоурутинов
	var chLocations = make(chan entity.Location)
	var chPrepareJobs = make(chan entity.JobForPrepare)
	var chGetJobs = make(chan entity.GotJob)
	var chJobsForSave = make(chan entity.JobForSave)
	//Признак занятости каналов
	var chGetListCount = make(chan byte, GetListC)
	var chPrepareJobCount= make(chan byte, PrepareJobC)
	var chGetJobCount= make(chan byte, GetJobC)

	//Ожидать
	var wGL  sync.WaitGroup
	var wPG  sync.WaitGroup
	var wGJ  sync.WaitGroup
	var wSG  sync.WaitGroup

	//Блок лист
	var blockList = block_list.NewBlockList()

	go rep.get_locations(chLocations)


	wGL.Add(GetListC)
	for i:=0; i<GetListC; i++{
		chGetListCount <- byte(1)
		go conveer_get_list(chLocations, chPrepareJobs, blockList, source, i, chGetListCount,  &wGL)
	}

	wPG.Add(PrepareJobC)
	for i:=0; i<PrepareJobC; i++{
		chPrepareJobCount<- byte(1)
		go conveer_prepare_job(chPrepareJobs, chGetJobs, blockList, source,  i, rep, chPrepareJobCount,  &wPG)
	}

	wGJ.Add(GetJobC)
	for i:=0; i<GetJobC; i++{
		chGetJobCount<- byte(1)
		go conveer_get_job(chGetJobs, chJobsForSave, source,  i, chGetJobCount,  &wGJ)
	}

	wSG.Add(SaveJobC)
	for i:=0; i<SaveJobC; i++{
		go conveer_save_job(chJobsForSave,source, i, rep, &wSG)
	}

	wGL.Wait()
	wPG.Wait()
	wGJ.Wait()
	wSG.Wait()
	defer wm.Done()
}

func (rep *Crawler)get_locations(chLocations chan entity.Location){
	loggingDB("get_locations",0,"Start coroutine")

	/*loc , err := rep.locations.FindByID(22)
	if err != nil {
		logging("get_locations",0,"Что то пошло не так")
		close(chLocations)
		return
	}

	chLocations <- *loc*/
	//TODO целевой
	locations, err :=rep.locations.All()
	if err != nil {
		logging("get_locations",0,"Что то пошло не так")
		close(chLocations)
		return
	}
	for _,location := range *locations {
		chLocations <- location
	}

	close(chLocations)
	logging("get_locations",0,"End coroutine")

}


func (rep *Crawler)compareWithIndex(arg *entity.CompareWithIndexIn) *entity.CompareResult {

	result := entity.NewCompareResult()

	// Компания в блок листе
	if strings.EqualFold(arg.CompanyName, "Jora Local"){
		result.ToBlockList()
		return result
	}

	//Пробуем найти компанию
	company, err := rep.company.FindCompanyByName(arg.CompanyName)
	if err != nil{
		result.ToAnyError()
		log.Info().Msg("compare_with_index: ошибко поиска в БД (FindCompanyByName)")
		log.Error().Msg(err.Error())
		return result

	}

	//Если не нашли компанию, то выходим. Поиск закончен
	if company == nil {
		return result
	}

	//Если нашли компанию, то отправим засуним ее в ответ
	result.CompanyId = *company

	//Пробуем найти работу
	job, err :=rep.job.FindJobBySourceOrTitle(arg.Id, result.Id, arg.Url, arg.TitleJob)
	if err != nil{
		result.ToAnyError()
		log.Info().Msg("compare_with_index: ошибко поиска в БД (FindJobBySourceOrTitle)")
		log.Error().Msg(err.Error())
		return result
	}

	//Если работа не найдена, то выходим
	if job == nil{
		return result
	}

	//Проверяем, совпадают URL или нет. Если совпадают, то ошибка
	if strings.EqualFold(job.Source, arg.Url){
		result.ToSameSource()
		return result
	}

	//Проверяем совпадение по компании и наименованию
	if strings.EqualFold(job.Title, arg.TitleJob){
		result.ToSameTitle()
		return result
	}

	log.Info().Msg("compare_with_index: Сюда мы не должны были дойти")
	result.ToAnyError()
	return result
}