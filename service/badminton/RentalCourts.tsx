import { Response } from '../../models/resp/base'
import { GetRentalCourts as Data } from '../../models/resp/rental-courts'
import { GetAuthRequestInit } from '../../data/auth/Auth'
import { GetRentalCourts as GetRentalCourtsApi } from '../../data/api/RentalCourts'

export async function GetRentalCourts(fromDate: Date, toDate: Date): Promise<Response<Data>> {
  const init = GetAuthRequestInit()
  return await GetRentalCourtsApi(fromDate, toDate, init)
}