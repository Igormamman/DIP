import { useGalleryPage } from "./hooks";
import React from "react";
import { CardGrid } from "../../Components/RaceCard/RaceCard";
import DatePicker from 'react-datepicker';
import { Container, Dropdown, Row, SplitButton, Col } from "react-bootstrap";
import "../../utils/styles.scss"

export const RaceGallery: React.FC = () => {

    const {
        races, startDate,
        endDate, onDateChange, loader, raceType, updateRaceType,
        userStore
    } = useGalleryPage()

    return (
        <Container>
            <Row>
                <Container id="layout-filterbar">
                    <Container className="margin-header white-space-nowrap" >
                        <Row>
                            <Col className="input margin">
                                <span style={{ verticalAlign: "middle" }}>Период дат:</span>
                                <DatePicker
                                    className="ui-inputtext ui-widget ui-corner-all"
                                    selected={startDate}
                                    onChange={onDateChange}
                                    startDate={startDate}
                                    endDate={endDate}
                                    fixedHeight
                                    showYearDropdown
                                    showMonthDropdown
                                    dropdownMode="select"
                                    selectsRange
                                />
                            </Col>
                            {userStore.isAuthorized &&
                                <Col className="input margin">
                                    <span style={{ verticalAlign: "middle" }}>По типу:</span>
                                    <SplitButton
                                        key={'primary'}
                                        id={`dropdown-split-variants-${'primary'}`}
                                        variant={'primary'.toLowerCase()}
                                        title={raceType[0]}
                                    >
                                        {raceType.slice(1).map((value, index) => <Dropdown.Item eventKey={index} onClick={() => updateRaceType(index + 1)}>{value}</Dropdown.Item>)}
                                    </SplitButton>
                                </Col>}
                        </Row>
                    </Container>
                </Container>
            </Row>
            <Row className="page-content">
                <CardGrid races={races} />
                {<div className="loader" ref={loader} />}
            </Row>
        </Container>
    )
};
