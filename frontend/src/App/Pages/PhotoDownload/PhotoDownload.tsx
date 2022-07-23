import React from "react";
import "../../utils/styles.scss"
import { useDownloadPage } from "./hooks";
import 'react-notifications-component/dist/theme.css'
import "animate.css"
import { Container, Row, Col, ProgressBar } from "react-bootstrap";
import { MyPagination } from "../../Components/Pagination";
import { PhotoGrid } from "../../Components/MasonryGrid/MasonryGrid";
import { PhotoDeletePreview } from "../../Components/Photo/PhotoDeletePreview";

export const DownloadPage: React.FC = () => {

    const {
        fileSelectedHandler,
        uploadFileHandler,
        access,
        selectedFileInfo,
        loading,
        fileCount,
        currentFileNum,
        handlePageClick,
        page,
        totalPages,
        UUIDs,
        raceAccess,
        deleteHandler,
    } = useDownloadPage()

    return (
        <>
            {(access && (raceAccess?.isMediaConfirmed == true || raceAccess?.isOrg == true)) &&
                <>
                    <Container >
                        <Row>
                            <Container id="layout-filterbar">
                                <Container className="margin-header white-space-nowrap" >
                                    <Row>
                                        <Col className="input margin  overflow-hidden">
                                            <h6> Выбор фото для загрузки </h6>
                                            <span>Название соревнования:
                                                <input id="getFile" type="file" multiple
                                                    onChange={(e) => fileSelectedHandler(e)} />
                                                <label style={{ marginLeft: 20 }} className="ng-tns-c9-2 ui-inputtext ui-widget ui-state-default ui-corner-all  " htmlFor="getFile">
                                                    {selectedFileInfo()}
                                                </label>
                                            </span>
                                        </Col>
                                    </Row>
                                </Container>
                            </Container>
                        </Row>
                        <Row className="page-content">
                            {!loading ?
                                <>
                                    <Col className="text-center" style={{ marginTop: 50 }}>
                                        <button className="btn btn-primary" onClick={uploadFileHandler}>
                                            Отправить изображения
                                        </button>
                                    </Col>
                                </> :
                                <>
                                    <Col style={{ margin: "50px auto" }} className="text-center">
                                        <ProgressBar animated visuallyHidden now={currentFileNum / (fileCount) * 100} label={`${currentFileNum}/${(fileCount)}`} style={{ height: 40, maxWidth: 800, margin: "auto" }} />
                                    </Col>
                                </>
                            }
                            <Container className="text-center" style={{ marginTop: "40px" }}>
                                {UUIDs.length>0 && <h1>Загруженные изображения:</h1>}
                                <MyPagination handlePage={handlePageClick} currentPage={page}
                                    totalPages={totalPages} />
                                {UUIDs.map((url, index) => (
                                    index % 2 == 0 &&
                                    <Row>
                                        <Col style={{ margin: "10px auto" }}>
                                            <PhotoDeletePreview url={`https://photo.marshalone.ru/api/photo/file/get?type=resized&UUID=${UUIDs[index]}`} deleteHandler={deleteHandler} />
                                        </Col>
                                        {UUIDs[index+1] &&
                                            <Col className={"justify-content-center"} style={{ margin: "10px auto" }}>
                                                <PhotoDeletePreview url={`https://photo.marshalone.ru/api/photo/file/get?type=resized&UUID=${UUIDs[index+1]}`} deleteHandler={deleteHandler} />
                                            </Col>
                                        }
                                    </Row>
                                )
                                )
                                }
                                <MyPagination handlePage={handlePageClick} currentPage={page}
                                    totalPages={totalPages} />
                            </Container>
                        </Row>
                    </Container>
                </>
            }
        </>)
};

