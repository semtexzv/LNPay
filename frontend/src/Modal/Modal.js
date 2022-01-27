import './Modal.css';
import React from 'react';

export function Modal({open, onDismiss, children}) {
    return (<div
        onClick={onDismiss}
        style={{display: open ? 'flex' : 'none'}} className={"modal"}>
        <div className={"modal-content"}>{children}</div>
    </div>);
}