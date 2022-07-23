import "../../utils/styles.scss"
import { useGalleryPage } from "./hooks";
import { Button } from "react-bootstrap";
import { MyPagination } from '../../Components/Pagination';
import React from "react";
import { PhotoGrid } from "../../Components/MasonryGrid/MasonryGrid";
import { RaceInfo } from "../../Components/RaceInfo/RaceInfo";
import { Container, Row, Col } from "react-bootstrap";
import { PhotoStore } from "../../Stores/PhotoStore";
import { observer, Observer } from "mobx-react-lite";


const RacePage: React.FC = () => {

    const {
        handlePageClick,
        gotoDownload,
        handleCompetitorNumber, undefinedPhotoClick,
        raceInfo, isActive,
        handleKeyDown, goToCart,
        raceAccess, photoStore, raceID, UUIDs
    } = useGalleryPage()


    const isActiveFunc = () => {
        return (
            <Observer>{() =>
                <Container onKeyDown={handleKeyDown}>
                    <Row>
                        <Container id="layout-filterbar">
                            <Container className="margin-header white-space-nowrap" >
                                <Row>
                                    <Col className="input margin" style={{ maxWidth: "fit-content" }}>
                                        <div className="form-switch d-flex align-items-center" style={{ height: "38px", width: "fit-content" }}>
                                            <span>Распознанные фото:</span>
                                            <input className="form-check-input" type="checkbox" id="flexSwitchCheckDefault" checked={photoStore[raceID].detected} onChange={undefinedPhotoClick} />
                                        </div>
                                    </Col>
                                    {photoStore[raceID].detected == true ?
                                        <>
                                            <Col className="input margin" style={{ width: "fit-content" }}>
                                                <div className="d-flex align-items-center" style={{ height: "38px" }}>
                                                    <span>Номер участника:</span>
                                                    <input className="ui-inputtext ui-widget ui-state-default ui-corner-all" style={{ verticalAlign: "middle" }} type="number" onChange={handleCompetitorNumber} />
                                                </div>
                                            </Col>
                                            {(raceAccess?.isMediaConfirmed == true || raceAccess?.isOrg == true) && <Col className="input margin">
                                                <Button variant="primary" onClick={gotoDownload}>Загрузить фото</Button>
                                            </Col>}
                                        </> :
                                        <>
                                            {(raceAccess?.isMediaConfirmed == true || raceAccess?.isOrg == true) && <Col className="input margin">
                                                <Button variant="primary" onClick={gotoDownload}>Загрузить фото</Button>
                                            </Col>}
                                        </>

                                    }
                                    <Col className="input margin">
                                        <Button variant="primary" onClick={goToCart}>К пользовательской галерее</Button>
                                    </Col>
                                </Row>
                            </Container>
                        </Container>
                    </Row>
                    <Row className="racepage-content">
                        <Container>
                            {raceInfo && <RaceInfo race={raceInfo} />}
                            <MyPagination handlePage={handlePageClick} currentPage={photoStore[raceID].currentPage}
                                totalPages={photoStore[raceID].totalPages} />
                            {UUIDs && UUIDs.length > 0 && <PhotoGrid uuids={UUIDs} />}
                            <MyPagination handlePage={handlePageClick} currentPage={photoStore[raceID].currentPage}
                                totalPages={photoStore[raceID].totalPages} />
                        </Container>
                    </Row>
                </Container>}
            </Observer>
        )
    }

    const isUnactiveFunc = () => {
        return (
            <>
                <Container style={{ textAlign: "center", marginTop: "60px" }}>
                    {raceInfo && <RaceInfo race={raceInfo} />}
                    <h1 style={{ margin: "200px auto", fontSize: "3em" }}>
                        Данная гонка еще не началась
                    </h1>
                </Container>
            </>
        )
    }

    return (
        <>
            {isActive == false && isUnactiveFunc()}
            {isActive == true && isActiveFunc()}
        </>
    )
};

export default observer(RacePage)
