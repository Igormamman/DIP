import React from "react";
import { Container, ProgressBar, Col, Row } from "react-bootstrap";
import { useCartPage } from "./hooks";



export const CartPage: React.FC = () => {

  const {
    items, DownloadButtonHandler, DeleteFromCart, loading, loadState
  } = useCartPage()

  const IsLoading = () => {
    return (
      <>
        {!loading ? items && items.map((item) => (
          <>
            <Col style={{ margin: "10px auto" }}>
              <img style={{ maxWidth: 240 }} src={`https://photo.marshalone.ru/api/photo/file/get?UUID=${item}&type=resized`} />
              <button className="btn btn-danger" style={{ marginLeft: 20 }} onClick={() => DeleteFromCart([item])}>Удалить</button>
            </Col>
          </>
        )) : items && <>
          <Col style={{ margin: "10px auto" }} className="text-center">
            <ProgressBar animated visuallyHidden now={loadState / (items?.length) * 100} label={`${loadState}/${(items?.length)}`} style={{ height: 40, maxWidth: 800, margin:"auto" }} />
          </Col>
        </>}
      </>
    )
  }

  return (
    <>
      <Container style={{ marginTop: "100px" }} className="text-center">
        <IsLoading />
        <Col>
          {items && items.length > 0 ? <button className="btn btn-primary" style={{ margin: "auto", alignSelf: "center" }}
            onClick={DownloadButtonHandler}>Скачать фото
          </button> : <h1 style={{ margin: "auto" }}>КОРЗИНА ПУСТА</h1>}
        </Col>
      </Container>
    </>)
};

