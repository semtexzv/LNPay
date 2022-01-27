import React, {useState, useCallback} from 'react';
import {ErrorRow, Field, MultiField, SubmitButton} from "./Fields";

export default function Complete({invoice, result}) {

    const [email, setEmail] = useState(null);

    const handleSend = useCallback(() => {
        console.log("Send email");
    }, [email]);

    return (<>
            <form onSubmit={() => {
            }}>
                <MultiField label={"Invoice"} disabled value={invoice.invoice}/>
                <Field label={"Amount"} disabled value={`${invoice.amount} ${invoice.currency} `}/>
                <Field label={"Payment id"} disabled value={result.charge_id}/>
                <Field label={"PreImage"} disabled value={result.preimage}/>
            </form>

            <form
                style={{marginTop: '1rem'}}
                onSubmit={handleSend}>
                <Field
                    label={"Email"}
                    value={email}
                    placeholder={'Enter email'}
                    onChange={e => setEmail(e.target.value)}
                    type='email'
                    autoComplete="true"
                />
                <SubmitButton disabled={!email}>Send me payment details</SubmitButton>
            </form>
        </>
    );

}