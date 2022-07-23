import React, { useEffect, useState, useCallback, useRef, RefObject } from "react";
import { useApi } from "../../../api/api"
import { useLocation, useNavigate } from "react-router-dom";
import { AppContext } from "../../AppContext"
import { Race } from "../../utils/types";
import { ServerRouts } from "../../utils/routs"
import moment from 'moment'
import { isSameDate } from "../../utils/common";


export const useGalleryPage = () => {

    const { api } = useApi()
    const navigate = useNavigate()
    const location = useLocation()
    const [races, setRaces] = useState<Race[]>([])
    const { userStore, calendarStore } = React.useContext(AppContext)
    // Flag to prevent simultaneous search operations 
    const [isLoading, setLoading] = useState(false)


    const loader = useRef(null);

    var setStart: Date
    var setEnd: Date
    var raceTypeList = ["Все гонки", "Медиа", "Участник"]
    if (calendarStore.raceType != "") {
        let tmp = raceTypeList;
        let index = raceTypeList.indexOf(calendarStore.raceType);
        [tmp[0], tmp[index]] = [tmp[index], tmp[0]];
        raceTypeList = tmp
    }
    if (!isSameDate(calendarStore.startDate, moment(new Date()).subtract(3, 'month').toDate())) {
        setStart = calendarStore.startDate
    } else {
        setStart = moment(new Date()).subtract(3, 'month').toDate()
    }
    if (!isSameDate(calendarStore.endDate, moment(new Date()).add(3, 'month').toDate())) {
        setEnd = calendarStore.endDate
    } else {
        setEnd = moment(new Date()).add(3, 'month').toDate()
    }
    // Race Class initialization 
    const [raceType, setRaceType] = useState(raceTypeList)
    const [isComp, setComp] = useState(raceTypeList[0] == "Участник" ? true : false)
    const [isOrg, setOrg] = useState(raceTypeList[0] == "Организатор" ? true : false)
    const [isMedia, setMedia] = useState(raceTypeList[0] == "Медиа" ? true : false)
    // Calendar params
    const [startDate, setStartDate] = useState(setStart);
    const [endDate, setEndDate] = useState(setEnd);

    const isOnScreen = useOnScreen(loader);

    const [pageLoaded, setPage] = useState(false)

    useEffect(() => {
        if (pageLoaded) {
            if (isOnScreen) {
                if (races.length > 0) {
                    window.scroll(0, window.scrollY * 0.98);
                    setStartDate((startDate) => moment(startDate).subtract(1, 'month').toDate())
                }
            }
        }
    }, [isOnScreen])

    // Calendar component handler    
    const onDateChange = (dates: any) => {
        const [start, end] = dates;
        setLoading(true)
        setStartDate(start);
        setRaces([])
        setEndDate(end);
        if (end != undefined && !isNaN(end.getTime())) {
            setLoading(false)
            loadRace(false)
        }
    };

    const resetDate = () => {
        setStartDate(moment(new Date()).subtract(3, 'month').toDate())
        setEndDate(moment(new Date()).add(3, 'month').toDate())
    }

    // function to load new race into raceList    
    const loadRace = async (check: boolean) => {
        if (isLoading) {
            return
        }
        await setLoading(true);
        var dateTo: string
        if (races != undefined && races.length > 0 && check) {
            dateTo = moment(races[races.length - 1].date).subtract(1, 'day').format('YYYY-MM-DD')
        } else {
            dateTo = moment(endDate).format('YYYY-MM-DD')
        }
        api.get(ServerRouts.Races, {
            withCredentials: true,
            params: {
                date_from: moment(startDate).format('YYYY-MM-DD'),
                date_to: dateTo,
            }
        })
            .then((response) => {
                let raceList: Race[]
                raceList = response.data
                raceList.reverse()
                if (raceList == undefined || raceList.length == 0) {
                    raceList = []
                }
                if (isComp) {
                    raceList = raceList.filter(e => e.isCompetitor == true)
                } else if (isMedia) {
                    raceList = raceList.filter(e => e.isMedia == true || e.isOrg == true)
                }
                setRaces(races => races.concat(raceList))
                setLoading(false)
            })
    }
    //"http://localhost:4005"
    //"https://user-dip.duckdns.org"

    // useEffect group   
    useEffect(() => {
        if (pageLoaded) {
            if (endDate != undefined && !isNaN(endDate.getTime())) {
                loadRace(true)
            }
        }
    }, [startDate, endDate]);

    useEffect(() => {
        if (pageLoaded) {
            setRaces([])
            loadRace(false)
            calendarStore.setRaceType(raceType[0])
        }
    }, [isMedia, isComp]);


    useEffect(() => {
        switch (raceType[0]) {
            case "Участник": {
                setMedia(false); setComp(true);
                break;
            }
            case "Медиа": {
                setMedia(true); setComp(false);
                break;
            }
            default: {
                setMedia(false); setComp(false);
                break;
            }
        }
    }, [raceType]);

    const updateRaceType = (index: number) => {
        let data = [...raceType];
        [data[0], data[index]] = [data[index], data[0]];
        setRaceType(data);
    }

    useEffect(() => {
        if (races.length > 0) {
            saveCalendarState()
        }
    }, [races]);

    useEffect(() => {
        if (calendarStore.raceList != undefined && calendarStore.raceList.length > 0) {
            setRaces(calendarStore.raceList)
        } else {
            loadRace(false)
        }
        setPage(true)
    }, []);

    const saveCalendarState = () => {
        calendarStore.setStartDate(startDate)
        calendarStore.setEndDate(endDate)
        calendarStore.setRaceList(races)
    }

    const gotoDownload = () => {
        navigate("../")
    }

    return {
        gotoDownload, races,
        startDate, onDateChange, endDate, loader, setComp, setMedia, setOrg, raceType, updateRaceType,
        resetDate, userStore
    }
}

export default function useOnScreen(ref: RefObject<HTMLElement>): boolean {
    const observerRef = useRef<IntersectionObserver>();
    const [isOnScreen, setIsOnScreen] = useState(false);

    const [width, setWidth] = useState<number>(window.innerWidth);

    function handleWindowSizeChange() {
        setWidth(window.innerWidth);
    }
    useEffect(() => {
        window.addEventListener('resize', handleWindowSizeChange);
        return () => {
            window.removeEventListener('resize', handleWindowSizeChange);
        }
    }, []);

    const isMobile = width <= 768;

    useEffect(() => {
        let rootMargin = isMobile ? '15%' : '1%' 
        observerRef.current = new IntersectionObserver(([entry]) =>
            setIsOnScreen(entry.isIntersecting), { rootMargin:  rootMargin}
        );
    }, []);

    useEffect(() => {
        if (observerRef.current != undefined && ref.current != undefined) {
            observerRef.current.observe(ref.current);
            return () => {
                if (observerRef.current != undefined) {
                    observerRef.current.disconnect();
                }
            };
        }
    }, [ref]);

    return isOnScreen;
}