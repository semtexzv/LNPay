import React, {useState, useMemo} from 'react';
import {CardElement, useElements, useStripe} from "@stripe/react-stripe-js";
import {CARD_OPTIONS, ErrorRow, Field, MultiField, SubmitButton} from "./Fields";
import {API_URL} from "../config";
import {decodeInvoice} from "../utils";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faStripe} from '@fortawesome/free-brands-svg-icons';


function CardSection() {
    return (
        <div className="FormRow" style={{flexDirection: 'column', marginBottom: '1px'}}>
            <CardElement options={CARD_OPTIONS}/>
            <div className={"Text"} style={{
                display: 'flex',
                alignItems: 'center'
            }}>Secured by &nbsp;<FontAwesomeIcon icon={faStripe} size={"2x"}/></div>
        </div>
    );
}

export default function CreateCharge({invoice, onFinish}) {
    const stripe = useStripe();
    const elements = useElements();
    const [processing, setProcessing] = useState(false);
    const [error, setError] = useState(null);
    const decoded = useMemo(() => decodeInvoice(invoice.invoice), [invoice]);

    console.log(invoice);

    const stripeTokenHandler = (token) => {
        console.log(token);
        fetch(`${API_URL}/api/v1/payment/${invoice.id}/charge`, {
            method: 'post',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({token: token['id']})
        })
            .then(res => {
                if (res.ok) {
                    return res.json();
                }
                throw res.json();
            })
            .then(res => {
                onFinish(res);
            })
            .catch(async err => {
                setError((await err)['error']);
            })
            .finally(() => setProcessing(false));
    };

    const handleSubmit = async (event) => {
        // We don't want to let default form submission happen here,
        // which would refresh the page.
        event.preventDefault();
        setProcessing(true);

        if (!stripe || !elements) {
            // Stripe.js has not yet loaded.
            // Make  sure to disable form submission until Stripe.js has loaded.
            return;
        }

        const card = elements.getElement(CardElement);
        const result = await stripe.createToken(card);

        if (result.error) {
            // Show error to your customer.
            console.log(result.error.message);
            setError(result.error.message);
            setProcessing(false);
        } else {
            // Send the token to your server.
            // This function does not exist yet; we will define it in the next step.
            stripeTokenHandler(result.token);
        }
    };

    return (
        <>
            <form onSubmit={handleSubmit}>
                <MultiField label={"Invoice"} disabled value={invoice.invoice}/>
                {decoded ? (
                    <Field
                        label={"Amount"}
                        value={`${decoded.valueSat} Satoshis (${decoded.valueSat / 100000000} BTC)`}
                        readOnly disabled
                    />) : null
                }
                <Field label={"Price"} disabled value={`~${invoice.amount} ${invoice.currency}`}/>
                <CardSection/>
                <ErrorRow error={error}/>
                <SubmitButton disabled={!stripe} processing={processing}>Confirm</SubmitButton>
            </form>
        </>
    );
}