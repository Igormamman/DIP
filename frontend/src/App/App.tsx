import React, { useEffect, useState } from "react";
import { Routes, Route } from 'react-router-dom';
import { useLocation } from 'react-router-dom';
import { AppRouts } from './utils/routs';
import { DownloadPage } from './Pages/PhotoDownload/PhotoDownload';
import { WatermarkPage } from "./Pages/WatermarkPage/WatermarkPage";
import { RaceGallery } from "./Pages/RaceGallery/RaceGallery";
import RacePage from "./Pages/RacePage/RacePage";
import { Header } from "./Components/M1Header/Header";
import { AppContext } from "./AppContext"
import { CartPage } from "./Pages/CartPage/CartPage";
import { ReactNotifications } from "react-notifications-component";
import { MyVerticallyCenteredModal } from "./Components/Photo/PhotoPreview";


export const App: React.FC = () => {

    const location = useLocation();
    let state = location.state as { backgroundLocation?: Location };

    return (
        <>
            <ReactNotifications />
            <Header />
            <Routes location={state?.backgroundLocation || location}>
                <Route path={AppRouts.Download} element={<DownloadPage />} />
                <Route path={`${AppRouts.Race}/:raceUID`} element={<RacePage />} />
                <Route path={AppRouts.FirstPage} element={<RaceGallery />} />
                <Route path={`${AppRouts.Watermark}/:UUID`} element={<MyVerticallyCenteredModal />} />
                <Route path={`${AppRouts.Download}/:raceUID`} element={<DownloadPage />} />
                <Route path={AppRouts.Cart} element={<CartPage />} />
            </Routes>
            {state?.backgroundLocation && (
                <Routes>
                    <Route path={`${AppRouts.Watermark}/:UUID`} element={<MyVerticallyCenteredModal />} />
                </Routes>
            )}
        </>
    );
}
