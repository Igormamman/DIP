import { makeAutoObservable } from "mobx";
import { ChangeEvent } from "react";

export class DownloadStore {

    private _currentFileNum: number

    private _totalFileCount: number

    private _isLoading: boolean

    constructor() {
        makeAutoObservable(this)
        this._currentFileNum = 0
        this._totalFileCount = 0
        this._isLoading = false
    }
   
    get currentFileNum(): number {
        return this._currentFileNum
    }

    setGalleryPage = (fileNum:number) => {
        this._currentFileNum = fileNum
    }

    get totalFileCount(): number {
        return this._totalFileCount
    }

    setTotalCount = (fileCount:number) => {
        this._totalFileCount = fileCount
    }

    get isLoading(): boolean {
        return this._isLoading
    }

    setLoading = (status:boolean) => {
        this._isLoading = status
    }

}