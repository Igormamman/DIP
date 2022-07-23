import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios";

const api_url =  'https://photo.marshalone.ru'
//'https://api-dip.duckdns.org'
// 'http://localhost:4000'
// 'http://photo.marshalone.ru'
let apiInstance: AxiosInstance

const createApi = ():AxiosInstance => {

    const api = axios.create({
        baseURL: api_url,
        timeout: 30000,
        withCredentials:true,
        validateStatus: function (status) {
            return status < 400; // Resolve only if the status code is less than 400
        }
    })

    const onSuccess = (response: AxiosResponse): AxiosResponse => {
        return response
    }

    const onFail = (error: AxiosError) => {
        return error
    }

    api.interceptors.response.use(onSuccess, onFail)

    return api

}

export const useApi = () => {
    if(!apiInstance)
        apiInstance = createApi()

    return {api: apiInstance}
}
