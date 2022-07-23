export enum AppRouts {
    FirstPage = '/',
    Download = '/download',
    Watermark = '/photo',
    Race = "/race",
    Cart = "/cart",
}


export enum ServerRouts {
    REGISTRATION = '/reg',
    LOGIN = 'api/users/login',
    LOGOUT = 'api/users/logout',
    Previews = 'api/photo/getPreviews',
    PhotoCount = 'api/photo/photo/count',
    MetaUpdate = 'api/photo/photo/task/meta',
    NewTask = 'api/photo/photo/task/new',
    Races = 'api/users/races',
    LoadPhoto = 'api/photo/file/upload',
    GetUserInfo='api/users/getUserInfo',
    GetRaceInfo = 'api/users/getRaceInfo',
    GetRaceAccess = 'api/users/getRaceAccess',
    DeletePhoto = 'api/photo/deletePhoto',
    GetFile = 'api/photo/file/get'
}

export enum M1Routs{
    GetRaceInfo = '/api/race'
} 
  