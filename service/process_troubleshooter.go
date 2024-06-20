package service

import (
	"xandrtools/client"
)

// send JSON an username and a password to a server
// get JSON with token
func getJobList(jobs []client.BatchSegmentUploadJob) []WebsiteBSUJ {
	joblist := make([]WebsiteBSUJ, len(jobs))
	for i, job := range jobs {
		joblist[i].BatchSegmentUploadJob = job
		joblist[i].BatchSegmentUploadJob.MatchRate = int(joblist[i].BatchSegmentUploadJob.NumValidUser * 100 / (joblist[i].BatchSegmentUploadJob.NumValidUser + joblist[i].BatchSegmentUploadJob.NumInvalidUser))
		if joblist[i].BatchSegmentUploadJob.MatchRate < 71 {
			joblist[i].BSUJerror.MatchRateErr = "Low match rate."
			joblist[i].JobErrors = append(joblist[i].JobErrors, joblist[i].BSUJerror.MatchRateErr)
		}
		if joblist[i].BatchSegmentUploadJob.ErrorLogLines != "" && joblist[i].BatchSegmentUploadJob.MatchRate < 71 {
			joblist[i].BSUJerror.ErrorLogLinesErr = "Remove invalid segments."
			joblist[i].JobErrors = append(joblist[i].JobErrors, joblist[i].BSUJerror.ErrorLogLinesErr)
		}
		if joblist[i].NumInvalidFormat > 0 {
			joblist[i].BSUJerror.NumInvalidFormatErr = "Fix invalid format."
			joblist[i].JobErrors = append(joblist[i].JobErrors, joblist[i].BSUJerror.NumInvalidFormatErr)
		}
		if joblist[i].NumUnauthSegment > 0 {
			joblist[i].NumUnauthSegmentErr = "Remove num_unauth_segment or verify that segment is active using apixandr.com/segment API call."
			joblist[i].JobErrors = append(joblist[i].JobErrors, joblist[i].NumUnauthSegmentErr)
		}
	}

	return joblist
}
