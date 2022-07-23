import React, { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom';
import { Card, Container, Row, Col } from "react-bootstrap";
import { useApi } from "../../../api/api";
import { Race } from '../../utils/types';
import { M1Routs, ServerRouts } from '../../utils/routs';
import { LazyLoadImage } from 'react-lazy-load-image-component';
import { AppContext } from '../../AppContext';
import axios from 'axios';
import "./styles.scss"
import moment from 'moment';
import 'moment/locale/ru';

interface Props {
    race:Race
}

export const RaceInfo: React.FC<Props> = (props) => {

    const [raceDate,setDate] = useState("")

    useEffect(()=>{
        moment.locale("ru");
        setDate(moment(props.race.date).format(('Do MMMM YYYY')))
    },[])

    return (
        <>
            {props.race &&
                <div style={{ textAlign:"center", width: "fit-content", margin:"0px auto 1.5em auto"}}>
                    <h1>{props.race.name}</h1>
                    <div className='date'>{raceDate}</div>
                    <div className='place'>{props.race.city}</div>
                </div>
            }
        </>
    )
}

