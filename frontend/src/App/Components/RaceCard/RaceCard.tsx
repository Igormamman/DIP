import React, { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Card, Container, Row, Col } from "react-bootstrap";
import { useApi } from "../../../api/api";
import "./styles.scss"
import { Race } from '../../utils/types';
import { ServerRouts } from '../../utils/routs';
import { LazyLoadImage } from 'react-lazy-load-image-component';
import { AppContext } from '../../AppContext';


interface Props {
  races: Race[]
}

export const RaceCard: React.FC<Race> = (race) => {

  const location = useLocation();



  return (
    <>
      <div style={{ width: "fit-content", margin: "auto" }}>
        <Link to={{
          pathname: `/race/${race.uid}`,
        }}>
          <Card style={{ width: "18rem", height: "25rem" }} >
            <LazyLoadImage
              src={"https://fget.marshalone.ru/files/race/uid/" + race.titlePicture}
              className="crop"
            />
            <Card.Body>
              <Card.Title>{race.name}</Card.Title>
              <Card.Subtitle className="mb-2 text-muted">Город:{race.city}</Card.Subtitle>
              <Card.Subtitle className="mb-2 text-muted">Дата проведения:{race.date}</Card.Subtitle>
            </Card.Body>
          </Card>
        </Link>
      </div>
    </>
  )
}



export const CardGrid: React.FC<Props> = (data: Props) => {

  const races = data.races

  const getColumnsForRow = () => {
    let items = races.map((race, _) => {
      return (
        <Col className="card-col">
          <RaceCard uid={race.uid} name={race.name} date={race.date} city={race.city} titlePicture={race.titlePicture} isMedia={false} isMediaConfirmed={false} isCompetitor={false} isOrg={false} ></RaceCard>
        </Col>
      );

    });
    return items;
  };

  return (
    <>
      {races && <div >
        <Container>
          <Row>
            {getColumnsForRow()}
          </Row>
        </Container>
      </div>}
    </>
  )
}
