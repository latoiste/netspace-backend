package db

import (
	"context"
	"time"
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
