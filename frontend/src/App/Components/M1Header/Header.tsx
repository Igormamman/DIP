import React, { useEffect, useState } from 'react';
import "./styles.scss"
import avatar from "../../utils/images/avatar.png"
import { AppContext } from '../../AppContext';
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom';
import { useApi } from '../../../api/api';
import { ServerRouts, AppRouts } from '../../utils/routs';
import { UserInfo } from '../../utils/types'
import { observer } from "mobx-react-lite"
import { Col, Container, Dropdown, Row } from 'react-bootstrap';
import axios from 'axios';

export const Header: React.FC = ({ }) => {
    const navigate = useNavigate()
    const { api } = useApi()
    const { userStore } = React.useContext(AppContext)
    const location = useLocation()

    const checkAuth = () => {
        if (!userStore.isAuthorized) {
            api.get(ServerRouts.GetUserInfo)
                .then((response) => {
                    let userInfo: UserInfo = {
                        cName: '',
                        cPhoto: ''
                    }
                    userInfo = response.data
                    if (response.status == 200) {
                        userStore.setAuthorized(true)
                        userStore.setUserName(userInfo.cName)
                        userStore.setUserPhoto(userInfo.cPhoto)
                    }
                })
        }
    }

    useEffect(() => {
        if (userStore.isAuthorized) {
            userStore.setAuthorized(true)
        } else {
            userStore.setAuthorized(false)
        }
    }, [userStore.isAuthorized])

    useEffect(() => {
        let image = new Image()
        image.src = avatar
        checkAuth()
    }, [])

    const onLoginClick = () => {
        if ('caches' in window) {
            caches.keys().then((names) => {
                // Delete all the cache files
                names.forEach(name => {
                    caches.delete(name);
                })
            });
        }
        window.location.assign("https://marshalone.ru/public/auth/");
    }

    const onLogoutClick = () => {
        api.get(ServerRouts.LOGOUT).then(
            () => { window.location.reload() },
            () => { window.location.reload() })
    }


    const goToM1 = () => {
        window.location.assign("https://marshalone.ru/");
    }

    const goToCart = () => {
        navigate(AppRouts.Cart)
    }

    const AvatarLine = observer(() => {
        if (userStore.isAuthorized) {
            return (<>
                <a className="inputgroup-addon overflow-hidden ng-star-inserted" onClick={onLogoutClick}>
                    {userStore.userPhoto != "" ? <><img className="avatar"
                        src={"https://fget.marshalone.ru/files/competitor/uid/" + userStore.userPhoto} />
                    </> : <><img className="avatar"
                        src={avatar} />
                    </>}
                    {userStore.userName}
                </a>
            </>)
        } else {
            return <><a className="inputgroup-addon overflow-hidden ng-star-inserted" onClick={onLoginClick}>Войти</a></>
        }
    }
    )

    const DropdownSelector = () => {

        type DropdownToggleProps = {
            children?: React.ReactNode;
            onClick?: (event: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {};
        };

        // The forwardRef is important!!
        // Dropdown needs access to the DOM node in order to position the Menu
        const DropdownToggle = React.forwardRef(
            (props: DropdownToggleProps, ref: React.Ref<HTMLAnchorElement>) => (
                <>
                    <a
                        href=""
                        ref={ref}
                        onClick={e => {
                            e.preventDefault();
                            if (props.onClick) props.onClick(e);
                        }}
                    >
                        {props.children}
                        <span style={{ paddingLeft: '5px' }}>&#x25bc;</span>
                    </a>
                </>
            )
        );

        type DropdownMenuProps = {
            children?: React.ReactNode;
            style?: React.CSSProperties;
            className?: string;
            labeledBy?: string;
        };

        // forwardRef again here!
        // Dropdown needs access to the DOM of the Menu to measure it
        const DropdownMenu = React.forwardRef(
            (props: DropdownMenuProps, ref: React.Ref<HTMLDivElement>) => {

                return (
                    <div
                        ref={ref}
                        style={props.style}
                        className={props.className}
                        aria-labelledby={props.labeledBy}
                    >
                        <ul className="list-unstyled">
                            {React.Children.toArray(props.children)}
                        </ul>
                    </div>
                );
            }
        );

        return (
            <Dropdown>
                <Dropdown.Toggle as={DropdownToggle} id="dropdown-custom-components" />
                <Dropdown.Menu as={DropdownMenu}>
                    <Dropdown.Item onClick={goToCart} eventKey="1">В корзину</Dropdown.Item>
                    <Dropdown.Item onClick={() => navigate(-1)} eventKey="2">Назад</Dropdown.Item>
                </Dropdown.Menu>
            </Dropdown>
        );
    };

    const headerGoBack = () => {
        if (location.pathname.startsWith("/race")) {
            navigate("..")
        } else if (location.pathname.startsWith("/download")) {
            var UUID = location.pathname.split('/download/')[1];
            navigate(`../race/${UUID}`)
        } else if (location.pathname.startsWith("/cart")) {
            navigate(-1)
        } else if (location.pathname.startsWith("/photo")) {
            navigate(`..`)
        } else {
            navigate(`..`)
        }
    }

    return (
        <>
            <Container fluid id="layout-topbar">
                <Row>
                    <Col className="d-md-none d-sm-none d-none d-lg-block"></Col>
                    <Col xs={10} className="content-header">
                        <Row className="align-items-center header">
                            <Col xs={4} className="white-space-nowrap text-overflow-ellipsis">
                                {(location.pathname != "/" && window.history.length > 1) &&
                                    <span className="logo end" onClick={headerGoBack}><i className="material-icons-round" style={{ verticalAlign: "middle", fontSize: 18 }}>arrow_back</i></span>
                                }
                            </Col>
                            <Col xs={4} className="center">
                                <span className="logo end" onClick={() => { navigate("") }}>photo.</span>
                                <span className="logo bold" onClick={goToM1}>marshal</span>
                                <span className="logo end" onClick={goToM1}>one</span>
                                <span className="logo end after" onClick={goToM1}>2</span>
                            </Col>
                            <Col xs={4}>
                                <Row className="white-space-nowrap overflow-hidden text-overflow-ellipsis text-right">
                                    <AvatarLine />
                                </Row>
                            </Col>
                        </Row>
                    </Col>
                    <Col className="d-md-none d-sm-none d-none d-lg-block col">
                    </Col>
                </Row>
            </Container >
        </>
    )
}






