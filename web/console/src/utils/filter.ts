import * as moment from '@/utils/moment'

export const dateRangeOptions: { value: string; label: string }[] = [
  {
    value: '0',
    label: '今天',
  },
  {
    value: '1',
    label: '过去7天',
  },
  {
    value: '2',
    label: '过去4周',
  },
  {
    value: '3',
    label: '过去3月',
  },
  {
    value: '4',
    label: '过去12月',
  },
  {
    value: '5',
    label: '本月至今',
  },
  {
    value: '6',
    label: '本季度至今',
  },
  {
    value: '7',
    label: '本年至今',
  },
  {
    value: '8',
    label: '所有时间',
  },
]

export const getRangeStartEndDates = (time_type: string): { start?: string; end?: string } => {
  const options = {}
  let start = ''
  if (time_type === '0') start = moment.getCurrentDate('YYYY-MM-DD 00:01')
  else if (time_type === '1') start = moment.getLastTimeAsDay(7, 'YYYY-MM-DD hh:mm')
  else if (time_type === '2') start = moment.getLastTimeAsWeek(4, 'YYYY-MM-DD hh:mm')
  else if (time_type === '3') start = moment.getLastTimeAsMonth(3, 'YYYY-MM-DD hh:mm')
  else if (time_type === '4') start = moment.getLastTimeAsMonth(12, 'YYYY-MM-DD hh:mm')
  else if (time_type === '5') start = moment.getCurrentMonth('YYYY-MM-DD hh:mm')
  else if (time_type === '6') start = moment.getCurrentQuarter('YYYY-MM-DD hh:mm')
  else if (time_type === '7') start = moment.getCurrentYeaer('YYYY-MM-DD hh:mm')
  else if (time_type === '8') start = '2022-01-01 00:00'

  if (start) {
    options.start = start
    options.end = moment.getCurrentDate('YYYY-MM-DD hh:mm')
  } else {
    options.start = null
    options.end = null
  }

  return options
}
