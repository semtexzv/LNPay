import React, {useState, useCallback, useMemo, useEffect} from 'react';
import {ErrorRow, Field, MultiField, SubmitButton,} from "./Fields";
import {API_URL} from "../config";
import {decodeInvoice, mapNetwork, removePrefix} from "../utils";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faQrcode, faSpinner, faTimes} from "@fortawesome/free-solid-svg-icons";
import {Modal} from "../Modal/Modal";


export function ReaderIcon({onInvoice, onError, onElemChange}) {
    const [scanning, setScanning] = useState(false);
    const [Element, setElement] = useState(false);

    const elem = useMemo(() => import('react-qr-scanner'), []);

    useEffect(() => {
        elem.then(e => {
            console.log(e);

            setElement(() => e.default);
        }).catch(() => {
            setScanning(false);
        });
    }, [elem, setElement]);

    const toggleScan = useCallback(() => {
        setScanning(!scanning);
    }, [scanning, setScanning]);

    const onScanResult = useCallback(data => {
        console.log("SCAN", data);
        if (data) {
            setScanning(false);
            if (data.startsWith('lightning:')) {
                onInvoice(removePrefix(data, "lightning:"));
            } else {
                onError("Not a valid lightning invoice");
            }
        }
    }, [setScanning, onInvoice, onError]);

    const onScanError = useCallback(err => {
        setScanning(false);
        console.log("SCAN", err);
        onError(new String(err));
    }, [setScanning, onError]);


    const scanIcon = scanning ? (Element ? faTimes : faSpinner) : faQrcode;

    const content = scanning && Element ? (
        <>
            <div style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
                <div className={"Text"}>Scan your invoice</div>
                <FontAwesomeIcon
                    style={{cursor: 'pointer'}}
                    onClick={toggleScan}
                    spin={scanning && !Element}
                    className="SocialIcon"
                    icon={faTimes}
                    size={'2x'}
                />

            </div>
            <Element
                style={{
                    width: '100%',
                    maxHeight: '100%'
                }}
                onClick={e => e.preventDefault()}
                onScan={onScanResult}
                onError={onScanError}
                facingMode={"rear"}
            />
        </>
    ) : null;

    return (
        <>
            <Modal
                open={scanning && Element}
                onDismiss={() => setScanning(false)}
            >
                {content}
            </Modal>
            <FontAwesomeIcon
                style={{cursor: 'pointer'}}
                onClick={toggleScan}
                spin={scanning && !Element}
                className="SocialIcon"
                icon={scanIcon}
            />
        </>
    );
}

export default function EnterInvoice({onFinish}) {
    const [invoice, _setInvoice] = useState("");
    const [invData, setInvData] = useState(null);
    const [error, setError] = useState(null);
    const [processing, setProcessing] = useState(false);
    const [scanner, setScanner] = useState(null);

    const changeInvoice = useCallback(invoice => {

        const decoded = decodeInvoice(invoice);
        if (!decoded) {
            setError("Invalid invoice");
        } else {
            setError(null);
        }
        console.log(decoded);
        _setInvoice(invoice);
        setInvData(invoice ? decoded : null);
    }, [_setInvoice]);

    const onsubmit = e => {
        e.preventDefault();
        setProcessing(true);
        fetch(`${API_URL}/api/v1/payment?invoice=${invoice}`, {method: 'post'})
            .then(res => {
                if (!res.ok) {
                    console.log("THrowing", res);
                    throw res.json();
                }
                return res.json();
            })
            .then(res => {
                onFinish(res);
            })
            .catch(async err => {
                setError((await err)['error']);
            })
            .finally(() => setProcessing(false));
    };

    return (
        <>
            {scanner}
            <form onSubmit={onsubmit}>
                <MultiField
                    label="Invoice"
                    id="invoice"
                    type="text"
                    placeholder="Lightning invoice"
                    required
                    value={invoice}
                    onChange={e => {
                        changeInvoice(e.target.value);
                    }}
                ><ReaderIcon onElemChange={setScanner} onInvoice={changeInvoice} onError={setError}/></MultiField>
                {invData ? (
                    <>
                        <Field label={"Network"} value={mapNetwork(invData.network)} readOnly disabled/>
                        <Field label={"Amount"}
                               value={`${invData.valueSat} Satoshis (${invData.valueSat / 100000000} BTC)`}
                               readOnly disabled/>
                        <Field label={"Description"} value={invData.shortDesc} readOnly disabled
                               placeholder={"No description"}/>

                    </>
                ) : null}
                <ErrorRow error={error}/>
                <SubmitButton disabled={error || !invoice} processing={processing}>Pay</SubmitButton>
            </form>
        </>
    );
}