import {useEffect, useState} from "react";
import {useApi} from "../../../api/api"
import {useNavigate} from "react-router-dom";
import queryString from "query-string"
import {PhotoMeta} from "../../utils/types";

export const useGalleryPage = () => {

    const {api} = useApi()
    const [url, setUrl] = useState("")
    const [meta, setMeta] = useState<PhotoMeta>()
    const navigate = useNavigate()
//api.get("http://localhost:4000/getPreviews",
    //api.get("https://api-dip.duckdns.org/getPreviews",


    useEffect(() => {
        const query = queryString.extract(window.location.href)
        loadMeta(query)
        setUrl(`https://api-dip.duckdns.org/photo/get?type=watermark&${query}`)
    }, []);

    const loadMeta = (query: string) => {
        api.get(`https://api-dip.duckdns.org/photo/meta?${query}`)
            .then((response) => {
                let photoMeta: PhotoMeta = response.data
                setMeta(photoMeta)
            })
    }

    const gotoGallery = () => {
        navigate("./..")
    }

    return {
        url,
        gotoGallery,
        meta
    }

}