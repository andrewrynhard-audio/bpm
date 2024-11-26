import React from "react";
import { OpenURL } from "../wailsjs/go/main/App";
import "./Modal.css";

interface ModalProps {
    show: boolean;
    title: string;
    message: string;
    url?: string;
    onClose: () => void;
}

const Modal: React.FC<ModalProps> = ({ show, title, message, url, onClose }) => {
    if (!show) return null;

    const handleOverlayClick = (e: React.MouseEvent<HTMLDivElement>) => {
        // Close the modal if the user clicks outside the modal content
        if ((e.target as HTMLElement).classList.contains("modal-overlay")) {
            onClose();
        }
    };

    const handleOpenURL = async () => {
        if (url) {
            await OpenURL(url);
        }
    };

    return (
        <div className="modal-overlay" onClick={handleOverlayClick}>
            <div className="modal-content">
                <h2>{title}</h2>
                <p>{message}</p>
                {url && (
                    <button onClick={handleOpenURL} className="update-link">
                        Click here to download
                    </button>
                )}
            </div>
        </div>
    );
};

export default Modal;
