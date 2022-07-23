export interface Race {
    uid: string
    name: string
    date: string
    city: string
    titlePicture: string
    isMedia: boolean
    isMediaConfirmed: boolean
    isCompetitor: boolean
    isOrg: boolean
}

export interface UserInfo {
    cName: string
    cPhoto: string
}

export interface Previews {
    count: number
    previewURL: string[]
}

export interface WidthHeight {
    width: number
    height: number
}

export interface PhotoMeta {
    UUID: string
    RUID: string
    height: number
    width: number
    competitors: string[]
}
