import { makeAutoObservable } from "mobx";
import moment from "moment";
import { ChangeEvent } from "react";
import { ThemeConsumer } from "react-bootstrap/esm/ThemeProvider";
import { Race } from "../utils/types";

export class CalendarStore {

    private _startDate: Date
    private _endDate: Date
    private _raceList: Race[]
    private _raceType: string

    constructor() {
        makeAutoObservable(this)
        this._startDate = moment(new Date()).subtract(3, 'month').toDate()
        this._endDate = moment(new Date()).add(3, 'month').toDate()
        this._raceList = []
        this._raceType = ""
    }

    get startDate(): Date {
        return this._startDate
    }

    get endDate(): Date {
        return this._endDate
    }

    get raceList(): Race[] {
        return this._raceList
    }

    get raceType(): string {
        return this._raceType
    }

    setRaceType = (raceType: string) => {
        if (["Все гонки", "Медиа", "Участник"].includes(raceType)){
            this._raceType = raceType
        } 
    }

    setStartDate = (startDate: Date) => {
        this._startDate = startDate
    }

    setEndDate = (endDate: Date) => {
        this._endDate = endDate
    }

    setRaceList = (raceList: Race[]) => {
        this._raceList = raceList
    }

    resetSearchStore = () => {
        this._startDate = moment(new Date()).subtract(3, 'month').toDate()
        this._endDate = moment(new Date()).add(3, 'month').toDate()
        this._raceList = []
    }

}