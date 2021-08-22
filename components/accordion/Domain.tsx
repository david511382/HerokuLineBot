

export interface AccordionData{
  summarys :string[]
  details?: AccordionData[]
}

export interface AccordionDataViewData  {
  start_id: number
  summarys :string[]
  details: AccordionDataViewData[]
}
