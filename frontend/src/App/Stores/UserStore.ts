import { makeAutoObservable } from "mobx";
import { ChangeEvent } from "react";

export class UserStore {

    private _userID: string
    private _userName: string
    private _photoID: string
    private _isAuthorized: boolean

    constructor() {
        makeAutoObservable(this)
        this._userID = ""
        this._userName = ""
        this._photoID = ""
        this._isAuthorized = false
    }

    get userID(): string {
        return this._userID
    }

    setUserID = (userID: string) => {
        this._userID = userID
    }

    get userName(): string {
        return this._userName
    }

    setUserName = (userName: string) => {
        this._userName = userName
    }

    get userPhoto(): string {
        return this._photoID
    }

    setUserPhoto = (photoID: string) => {
        this._photoID = photoID
    }

    get isAuthorized(): boolean {
        return this._isAuthorized
    }

    setAuthorized = (status: boolean) => {
        this._isAuthorized = status
    }

    resetUserStore = () => {
        this._isAuthorized = false
        this._userName = ""
    }

}