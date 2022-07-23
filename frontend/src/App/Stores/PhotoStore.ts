import { makeAutoObservable } from "mobx";

export class PhotoStore {

    private _UUIDs: string[]
    private _totalPages: number
    private _currentPage: number
    private _detected: boolean

    constructor() {
        makeAutoObservable(this)
        this._UUIDs = []
        this._totalPages = 1
        this._currentPage = 1
        this._detected = true
    }

    get detected(): boolean {
        return this._detected
    }

    setDetected = (state: boolean) => {
        this._detected = state
    }

    get UUIDs(): string[] {
        return this._UUIDs
    }

    setUUIDs = (uuids: string[]) => {
        this._UUIDs = uuids
    }

    get totalPages(): number {
        return this._totalPages
    }

    setTotalPages = (count: number) => {
        this._totalPages = count
    }

    get currentPage(): number {
        return this._currentPage
    }

    setCurrentPage = (page: number) => {
        this._currentPage = page
    }

    resetState = () => {
        this.setUUIDs([])
        this.setTotalPages(1)
        this.setCurrentPage(1)
        this.setDetected(true)
    }

}