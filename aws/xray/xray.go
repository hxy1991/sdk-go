package awsxray

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func GetSegmentAndTraceId(ctx context.Context) (string, string) {
	defer func() {
		if r := recover(); r != nil {
			// ignore
		}
	}()

	sCtx, seg := xray.BeginSubsegment(ctx, "zap-x-ray-log")
	seg.Close(nil)

	traceId := seg.TraceID

	if traceId == "" {
		traceId = xray.TraceID(sCtx)
	}

	if traceId == "" {
		traceId = xray.TraceID(ctx)
	}

	segmentId := seg.ParentSegment.ID

	seg.ParentSegment.RemoveSubsegment(seg)

	return segmentId, traceId
}
