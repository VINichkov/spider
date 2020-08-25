package crawler

import (
	"spider/block_list"
	"spider/models/entity"
	"spider/source"
	"strconv"
	"sync"
)


func conveer_get_list(chLocations chan entity.Location, job_for_prepare chan entity.JobForPrepare,
	block *block_list.BlockList, source source.Source, thread int, conv chan byte, wGl *sync.WaitGroup){
	loggingDB("conveer_get_list", thread,source.Name() +": Start coroutine")

	for location := range chLocations {
		getMainPage(location, job_for_prepare, block, source, thread)
	}

	<-conv
	if len(conv) == 0 {
		close(job_for_prepare)
	}

	loggingDB("conveer_get_list", thread, source.Name() +": End conv = " + strconv.Itoa(len(conv)))
	defer wGl.Done()

}

func conveer_prepare_job(chPrepareJobs chan entity.JobForPrepare, chGetJobs chan entity.GotJob,
	block *block_list.BlockList, source source.Source,  thread int, rep *Crawler, conv chan byte,
	wPJ *sync.WaitGroup){
	loggingDB("conveer_prepare_job", thread,source.Name() +": Start coroutine")

	for preapreJob := range chPrepareJobs {
		getList(preapreJob, chGetJobs, block, source,  thread, rep)
	}

	<-conv
	if len(conv) == 0 {
		close(chGetJobs)
	}

	loggingDB("conveer_prepare_job", thread,source.Name() +": End")
	defer wPJ.Done()


}

func conveer_get_job(chGetJobs chan entity.GotJob, chJobsForSave chan entity.JobForSave,
	source source.Source, thread int, conv chan byte, wGJ *sync.WaitGroup) {
	loggingDB("conveer_get_job", thread,source.Name() +": Start coroutine")

	for got_job := range chGetJobs {
		get_job(got_job, chJobsForSave, source, thread)
	}

	<-conv
	if len(conv) == 0 {
		close(chJobsForSave)
	}

	loggingDB("conveer_get_job", thread,source.Name() +": End")
	defer wGJ.Done()

}

func conveer_save_job(chJobsForSave chan entity.JobForSave, source source.Source,thread int,
	rep *Crawler, wSJ *sync.WaitGroup){
	loggingDB("conveer_save_job", thread,source.Name() +": Start coroutine")

	for jobForSave := range chJobsForSave {
		save_job(jobForSave,  thread, rep)
	}

	loggingDB("conveer_save_job", thread,source.Name() +": End")
	defer wSJ.Done()
}