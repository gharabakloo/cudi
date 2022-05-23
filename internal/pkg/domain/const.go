package domain

import "time"

const (
	Colon   = ":"
	AllTags = "*"

	Reference      = "reference"
	TypeSeparately = "separately"
	TypeTogether   = "together"
)

const (
	Year   = 365 * Day
	Month  = 30 * Day
	Week   = 7 * Day
	Day    = 24 * Hour
	Hour   = 60 * Minute
	Minute = 60 * Second
	Second = time.Second
)

const (
	YearSign   = "y"
	MonthSign  = "m"
	WeekSign   = "w"
	DaySign    = "d"
	HourSign   = "h"
	MinuteSign = "min"
	SecondSign = "s"
)
