package client

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalTime(t *testing.T) {
	var bsResponse BatchSegmentResponse

	if err := json.Unmarshal([]byte(text), &bsResponse); err != nil {
		t.Fatal(err)
	}

	ct := time.Time(bsResponse.Response.BatchSegmentUploadJob[0].CompletedTime)
	t.Log("completed_time:", ct.Format(time.RFC3339))
}

const text = `{
	"response": {
		"batch_segment_upload_job": [
			{
				"completed_time": "2024-06-03 21:00:00",
				"created_on": "2024-06-03 21:09:33",
				"error_code": null,
				"error_log_lines": "\n\nnum_unauth_segment-4013681496264948522;5013:0,5014:1550",
				"id": -6170722753141248786,
				"is_beam_file": false,
				"job_id": "e3cba40ed695fb88e60cf31dc40902c19f13ddbf",
				"last_modified": "2024-06-03 21:09:34",
				"member_id": 823,
				"num_inactive_segment": 0,
				"num_invalid_format": 0,
				"num_invalid_segment": 0,
				"num_invalid_timestamp": 0,
				"num_invalid_user": 0,
				"num_other_error": 0,
				"num_past_expiration": 0,
				"num_unauth_segment": 1,
				"num_valid": 200000,
				"num_valid_user": 100000,
				"percent_complete": 100,
				"phase": "completed",
				"segment_log_lines": "\n5010:100000\n5011:50000\n5012:50000",
				"start_time": "2024-06-03 21:09:33",
				"time_to_process": 60000001065,
				"uploaded_time": "2024-06-03 21:09:33",
				"validated_time": "2024-06-03 21:09:36"
			}
		],
		"count": 1,
		"dbg_info": {
			"instance": "authentication-api-production-8664bd4765-btqsz",
			"time": 0,
			"start_time": "2024-06-03T17:33:38.497553884-04:00",
			"version": "0.0.0",
			"trace_id": "d00284f4d8a9fc9a8f8f",
			"warnings": null
		},
		"num_elements": 0,
		"start_element": 0,
		"status": "OK"
	}
}
`
