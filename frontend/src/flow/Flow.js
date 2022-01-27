import EnterInvoice from "./EnterInvoice";
import React, {useState} from 'react';
import CreateCharge from "./CreateCharge";
import Complete from "./Complete";

const FLOW_ENTER_INVOICE = 0;
const FLOW_CREATE_CHARGE = 1;
const FLOW_COMPLETE = 2;

export const Flow = () => {
    const [stage, setStage] = useState(0);
    const [invoice, setInvoice] = useState(null);
    const [result, setResult] = useState(null);

    const onInvoice = invoice => {
        console.log("Paying invoice", invoice);
        setInvoice(invoice);
        setStage(FLOW_CREATE_CHARGE);
    };

    const onComplete = result => {
        setResult(result);
        setStage(FLOW_COMPLETE);
    };

    if (stage === FLOW_ENTER_INVOICE) {
        return (<EnterInvoice onFinish={onInvoice}/>);
    }

    if (stage === FLOW_CREATE_CHARGE) {
        return (<CreateCharge invoice={invoice} onFinish={onComplete}/>);
    }
    if (stage == FLOW_COMPLETE) {
        return (<Complete invoice={invoice} result={result} />);
    }
    return (<div> Not yet</div>);
};