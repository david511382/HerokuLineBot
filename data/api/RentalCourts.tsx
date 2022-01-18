import { GetBackendUrl } from '../env/Http'
import type { GetRentalCourts as Response } from '../../models/resp/rental-courts'

export async function GetRentalCourts(fromDate : Date, toDate : Date, init?: RequestInit | undefined) :Promise<Response> {
  const url = `${GetBackendUrl()}/api/badminton/rental-courts`
  const urlParams = new URLSearchParams();
  urlParams.set("from_date", fromDate.toISOString())
  urlParams.set("to_date", toDate.toISOString())
  return await fetch(
    `${url}?${urlParams.toString()}`,
    init,
  ).then((response) => {
    if (!response.ok) {
      throw new Error(response.statusText)
    }
    return response.json()
  })
}