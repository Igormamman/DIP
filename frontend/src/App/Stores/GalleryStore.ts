import { makeAutoObservable } from "mobx";
import { ChangeEvent } from "react";

export class GalleryStore {

    private _galleryPage: number


    constructor() {
        makeAutoObservable(this)
        this._galleryPage = 0
    }
    get galleryPage(): number {
        return this._galleryPage
    }

    setGalleryPage = (page:number) => {
        this._galleryPage = page
    }

}