import { DownloadStore } from "./DownloadStore"
import { GalleryStore } from "./GalleryStore"
import { UserStore } from "./UserStore"
import { CalendarStore } from "./CalendarStore"
import { PhotoStore } from "./PhotoStore"

export class MainStore {

    userStore: UserStore
    galleryStore: GalleryStore
    downloadStore: DownloadStore
    calendarStore: CalendarStore
    photoStore: {[raceUID:string]:PhotoStore}

    constructor() {
        this.userStore = new UserStore()
        this.galleryStore = new GalleryStore()
        this.downloadStore = new DownloadStore()
        this.calendarStore = new CalendarStore()
        this.photoStore = {}
    }
}