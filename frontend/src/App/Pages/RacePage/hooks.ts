import React, { useEffect, useState } from "react";
import { useApi } from "../../../api/api"
import { useNavigate, useParams } from "react-router-dom";
import { Previews, Race } from "../../utils/types";
import { AppRouts, ServerRouts } from "../../utils/routs"
import { debounce } from "../../utils/common";
import moment from "moment";
import { AppContext } from "../../AppContext";
import { PhotoStore } from "../../Stores/PhotoStore";

export const useGalleryPage = () => {

    const { photoStore } = React.useContext(AppContext)
    const { raceUID } = useParams();
    let raceID: string
    if (raceUID != undefined) {
        raceID = raceUID
    } else {
        raceID = ""
    }
    if (photoStore[raceID] == undefined) {
        photoStore[raceID] = new PhotoStore()
        photoStore[raceID].resetState()
    }

    const { api } = useApi()
    const [inputcompetitor, setInputCompetitor] = useState("")
    const [competitor, setCompetitor] = useState("")
    const [UUIDs, setUUIDs] = useState<string[]>([])
    const navigate = useNavigate()

    const [pageLoaded, setLoaded] = useState(false)
    const [raceInfo, setRaceInfo] = useState<Race>()
    const [raceAccess, setRaceAccess] = useState<Race>()
    const [isActive, setActive] = useState<boolean>()

    const getRaceInfo = async () => {
        api.get(ServerRouts.GetRaceInfo, {
            params: {
                ruid: raceUID
            }
        }).then((response) => {
            let data: Race
            if (response.data != undefined) {
                data = response.data
                if (data != null) {
                    setRaceInfo(data);
                    if (moment().diff(data.date, 'days') >= 0) {
                        setActive(true)
                    } else {
                        setActive(false)
                    }
                }
            }
        })
    }

    useEffect(()=>{
        photoStore[raceID].setUUIDs(UUIDs)
    },[UUIDs])

    const getRaceAccess = async () => {
        api.get(ServerRouts.GetRaceAccess, {
            params: {
                ruid: raceUID
            }
        }).then((response) => {
            let data: Race
            if (response.data != undefined) {
                data = response.data
                if (data != null) {
                    setRaceAccess(data);
                }
            }
        })
    }

    const handlePageClick = (event: any) => {
        photoStore[raceID].setCurrentPage(event);
    };

    const handleCompetitorNumber = (event: any) => {
        if (event.target.value) {
            setInputCompetitor(event.target.value)
            setCompetitor(event.target.value)
        } else {
            setInputCompetitor("")
            setCompetitor(event.target.value)
        }
    }

    useEffect(() => {
        getRaceInfo()
        getRaceAccess()
    }, [])

    useEffect(() => {
        getPreviewUrls()
        setLoaded(true)
    }, []);

    useEffect(() => {
        if (pageLoaded) {
            if (photoStore[raceID].currentPage == 1) {
                getPreviewUrls()
            } else {
                photoStore[raceID].setCurrentPage(1)
            }
        }
    }, [photoStore[raceID].detected])

    useEffect(() => {
        if (pageLoaded) {
            debouncedGetPreviewUrls()
        }
    }, [photoStore[raceID].currentPage, competitor]);

    //  api.get("https://api-dip.duckdns.org/photo/count"
    //  api.get("http://localhost:4000/photo/count"
    const getPreviewUrls = () => {

        api.get(ServerRouts.Previews, {
            params: {
                competitor: competitor,
                raceUID: raceUID,
                limit: 48,
                detected: photoStore[raceID].detected,
                offset: (photoStore[raceID].currentPage - 1) * 48,
            }
        })
            .then((response) => {
                let data: Previews
                if (response.data != undefined) {
                    data = response.data
                } else {
                    data = { count: 0, previewURL: [] }
                    photoStore[raceID].setTotalPages(1)
                    photoStore[raceID].setCurrentPage(1)
                    setUUIDs([])
                    return
                }
                setUUIDs(data.previewURL)
                if (data.count % 48 == 0) {
                    photoStore[raceID].setTotalPages(Math.floor(data.count / 48))
                } else {
                    photoStore[raceID].setTotalPages(Math.floor(data.count / 48) + 1)
                }
            })
    }

    const debouncedGetPreviewUrls = debounce(getPreviewUrls, 500);
    // api.get("https://api-dip.duckdns.org/getPreviews",
    //  api.get("http://localhost:4000/getPreviews",


    const undefinedPhotoClick = () => {
        setCompetitor("")
        if (photoStore[raceID].detected != false) {
            photoStore[raceID].setDetected(false)
        } else {
            photoStore[raceID].setDetected(true)
        }
    }

    const handleKeyDown = (event: React.KeyboardEvent<HTMLDivElement>) => {
        if (event.key === 'Enter') {
            photoStore[raceID].setCurrentPage(1)
            setCompetitor(inputcompetitor)
        }
    }

    const handleSearch = () => {
        photoStore[raceID].setCurrentPage(1)
        setCompetitor(inputcompetitor)
        //     setRace(selectedRace)
    }

    const goToCart = () => {
        navigate(AppRouts.Cart)
    }

    //api.get("http://localhost:4000/getPreviews",
    //api.get("https://api-dip.duckdns.org/getPreviews",

    const gotoDownload = () => {
        navigate(`../download/${raceUID}`)
    }

    return {
        handlePageClick, photoStore,
        handleSearch, gotoDownload,
        handleCompetitorNumber, undefinedPhotoClick, raceInfo, isActive,
        handleKeyDown, raceAccess, goToCart, raceID, UUIDs
    }

}
