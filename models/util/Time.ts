export const FORMATE_DATE = "yyyy-MM-dd"
export const FORMATE_RFC3339 = "yyyy-MM-ddThh:mm:ss.fffTzzz"

export enum TimeType {
	YEAR,
	MONTH,
	DATE,
}

export interface TimeInit {
	year: number,
	month?: number,
	date?: number,
	hours?: number,
	minutes?: number,
	seconds?: number,
	ms?: number,
}

export class Time extends Date {
	static Weekday(day : number) : string{
		switch (day){
			case 0: 
				return "日"
			case 1:
				return "一"
			case 2:
				return "二"
			case 3:
				return "三"
			case 4:
				return "四"
			case 5:
				return "五"
			case 6:
				return "六"
			default:
				return this.Weekday(day % 7)
		}
	}

	static Weekdays(startFrom? : number) : string[]{
		const result : string[] = []

		if (!startFrom)
			startFrom = 1
		
		for (let i = 0; i < 7; i++){
			let day = startFrom + i
			result.push(this.Weekday(day))
		}	

		return result
	}	

	
	static Slice(
		from : Time, beforeTime : Time,
		nextTime : (t : Time) => Time,
		doF : (runTime : Time, next : Time) => boolean,
	) {
		let runTime = from
		for (let dur = beforeTime.Sub(runTime); dur > 0; dur = beforeTime.Sub(runTime)) {
			let next = nextTime(runTime)
			if (next.valueOf() > beforeTime.valueOf()) {
				next = beforeTime
			}

			if (!doF(runTime, next)) {
				break
			}

			runTime = next
		}
	}


	format : string

	constructor (
		value?: number | string | Date | TimeInit,
		format? :string,
	){
		if (value){
			if (instanceOfTimeInit(value)){
				super(
					(value.year)?value.year:1,
					(value.month)?value.month - 1:1,
					(value.date)?value.date:1,
					(value.hours)?value.hours:0,
					(value.minutes)?value.minutes:0,
					(value.seconds)?value.seconds:0,
					(value.ms)?value.ms:0,
				)	

				this.format = (format)?
					format:
					(value.hours)?
						FORMATE_RFC3339:
						FORMATE_DATE
				return
			}else
				super(value)
		}
		else
			super()
 	
		this.format = (format)?
			format:
			FORMATE_RFC3339
	}

	Format(format? :string) : string {
		if (!format){
			format = this.format
		}
		
		let result : string = ""
		for (let i = 0; i < format.length; i ++){
			const c = format[i]
			let count = 1
			for (;i+1 < format.length;){
				const nextI = i + 1
				const next = format[nextI]
				if (next !== c)
					break

				count++
				i = nextI
			}

			let appendValue = ""
			switch (c){
				case 'y':
					appendValue = this.getFullYear().toString()
					break
				case 'M':
					appendValue = (this.getMonth()).toString()
					break
				case 'd':
					appendValue = this.getDate().toString()
					break
				case 'h':
					appendValue = this.getHours().toString()
					break
				case 'm':
					appendValue = this.getMinutes().toString()
					break
				case 's':
					appendValue = this.getSeconds().toString()
					break
				case 'f':
					appendValue = this.getMilliseconds().toString()
					break
				case 'w':
					appendValue = Time.Weekday(this.getDay())
					count = appendValue.length
					break
				case 'z':
					const offset = this.getTimezoneOffset()
					if (offset === 0){
						appendValue = "Z"
					}else {
						const signStr = (offset < 0)?
							"+":
							"-"

						const absOffset = Math.abs(offset)
						const absTimezoneOffsetHour = Math.abs(absOffset / 60)
						const absTimezoneOffsetHourStr = setStrLen(absTimezoneOffsetHour.toString(), 2)
						
						const absTimezoneOffsetMin = absOffset % 60
						const absTimezoneOffsetMinStr = setStrLen(absTimezoneOffsetMin.toString(), 2)

						appendValue = `${signStr}${absTimezoneOffsetHourStr}:${absTimezoneOffsetMinStr}`
					}
					count = appendValue.length
					break
				default:
					appendValue = c
					break
			}

			appendValue = setStrLen(appendValue, count)

			result += appendValue
		}
		
		return result 
	}

	Of(timeType : TimeType) : Time{
		const init : TimeInit = {
			year : this.getFullYear(),
		}
		switch (timeType){
			case TimeType.DATE:
				init.date = this.getDate()
			case TimeType.MONTH:
				init.month = this.getMonth()
		}
		return new Time(init)
	}

	getMonth() : number{
		return super.getMonth()+1
	}

	setMonth(month: number, date?: number | undefined): number{
		month = month - 1
		if (date)
			return super.setMonth(month, date)
		return super.setMonth(month)
	}

	Next(timeType : TimeType, count :number) : Time{
		let result = new Time(this.valueOf())
		switch (timeType) {
			case TimeType.DATE:
				result.setDate(this.getDate() + count)
				break
			case TimeType.MONTH:
				result.setMonth(this.getMonth() + count)
				break
			case TimeType.YEAR:
				result.setFullYear(this.getFullYear() + count)
				break
		}

		return result
	}


	Sub(date : Time | Date) : number {
		return this.valueOf() - date.valueOf()
	}

	Equal(target : Date) : boolean {
		return this.getTime() === target.getTime()
	}
}	

function instanceOfTimeInit(object: any): object is TimeInit {
	return (object as TimeInit).year !== undefined
}

function setStrLen(currentValue : string, wantLen : number) : string{
	const len = currentValue.length
	if (len > wantLen){
		const from = len - wantLen
		currentValue = currentValue.substr(from)
	}else if (len < wantLen) {
		const amount = wantLen - len
		currentValue = '0'.repeat(amount) + currentValue
	}
	
	return currentValue
}