package db

import (
	"context"
	"time"

	"github.com/latoiste/netspace/api"
)

func (r *Repository) CheckInAnalytics(
	yesterdayLocalStart time.Time,
	todayLocalStart time.Time,
	todayLocalCurrent time.Time,
	ctx context.Context,
) (today int, yesterday int, err error) {
	today, err = r.TotalCheckInRange(
		todayLocalStart.UTC(),
		todayLocalCurrent.UTC(),
		ctx,
	)

	if err != nil {
		return 0, 0, err
	}

	yesterday, err = r.TotalCheckInRange(
		yesterdayLocalStart.UTC(),
		todayLocalStart.UTC(),
		ctx,
	)

	if err != nil {
		return 0, 0, err
	}

	return today, yesterday, nil
}

func (r *Repository) ActiveUsersAnalytics(
	yesterdayLocalStart time.Time,
	todayLocalStart time.Time,
	todayLocalCurrent time.Time,
	ctx context.Context,
) (today int, yesterday int, err error) {
	today, err = r.TotalActiveUsersRange(
		todayLocalStart.UTC(),
		todayLocalCurrent.UTC(),
		ctx,
	)
	if err != nil {
		return 0, 0, err
	}

	yesterday, err = r.TotalActiveUsersRange(
		yesterdayLocalStart.UTC(),
		todayLocalStart.UTC(),
		ctx,
	)
	if err != nil {
		return 0, 0, err
	}

	return today, yesterday, nil
}

func (r *Repository) MessagesAnalytics(
	yesterdayLocalStart time.Time,
	todayLocalStart time.Time,
	todayLocalCurrent time.Time,
	ctx context.Context,
) (today int, yesterday int, err error) {
	today, err = r.TotalMessagesRange(
		todayLocalStart.UTC(),
		todayLocalCurrent.UTC(),
		ctx,
	)
	if err != nil {
		return 0, 0, err
	}

	yesterday, err = r.TotalMessagesRange(
		yesterdayLocalStart.UTC(),
		todayLocalStart.UTC(),
		ctx,
	)
	if err != nil {
		return 0, 0, err
	}

	return today, yesterday, nil
}

func (r *Repository) HourlyCheckInAnalytics(bucket int, now time.Time, ctx context.Context) ([]api.HourlyCheckInDTO, error) {
	results := make([]api.HourlyCheckInDTO, bucket)

	currentHour := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(), 0, 0, 0,
		now.Location(),
	)

	for i := range bucket {
		start := currentHour.Add(time.Duration(-i) * time.Hour)
		end := start.Add(time.Hour)

		count, err := r.TotalCheckInRange(
			start.UTC(),
			end.UTC(),
			ctx,
		)
		if err != nil {
			return nil, err
		}

		dto := api.HourlyCheckInDTO{
			Label: start.Format("15:04"),
			Value: count,
		}
		results[bucket-1-i] = dto
	}

	return results, nil
}
