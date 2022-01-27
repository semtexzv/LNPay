import React from 'react';
import {CardElement} from "@stripe/react-stripe-js";
import TextareaAutosize from 'react-textarea-autosize';

export const CARD_OPTIONS = {
    iconStyle: 'solid',
    /*
    style: {
        base: {
            iconColor: '#c4f0ff',
            color: '#fff',
            fontWeight: 500,
            fontFamily: 'Roboto, Open Sans, Segoe UI, sans-serif',
            fontSize: '16px',
            fontSmoothing: 'antialiased',
            ':-webkit-autofill': {
                color: '#fce883',
            },
            '::placeholder': {
                color: '#87bbfd',
            },
        },
        invalid: {
            iconColor: '#ffc7ee',
            color: '#ffc7ee',
        },
    },*/
};

export const Field = ({
                          label,
                          id,
                          type,
                          placeholder,
                          required,
                          autoComplete,
                          value,
                          onChange,
                          ...props
                      }) => (
    <div className="FormRow">
        <label htmlFor={id} className="FormRowLabel">
            {label}
        </label>
        <input
            className="FormRowInput"
            id={id}
            type={type}
            placeholder={placeholder}
            required={required}
            autoComplete={autoComplete}
            value={value}
            onChange={onChange}
            {...props}
        />
    </div>
);

export const ErrorRow = ({error}) => {
    return (
        <div className="FormRow">
            <label className="FormRowInput">
                {error}
            </label>
        </div>
    );
};

export const MultiField = ({
                               label,
                               id,
                               type,
                               placeholder,
                               required,
                               autoComplete,
                               value,
                               onChange,
                               children,
                               ...props
                           }) => (
    <div className="FormRow">
        <label htmlFor={id} className="FormRowLabel">
            {label}
        </label>
        <TextareaAutosize
            className="FormRowInput"
            id={id}
            type={type}
            placeholder={placeholder}
            required={required}
            autoComplete={autoComplete}
            value={value}
            onChange={onChange}
            {...props}
        />{children}
    </div>
);


export const SubmitButton = ({processing, error, children, disabled}) => (
    <button
        className={`SubmitButton ${error ? 'SubmitButton--error' : ''}`}
        type="submit"
        disabled={processing || disabled}
    >
        {processing ? 'Processing...' : children}
    </button>
);

export const CardField = ({onChange}) => (
    <div className="FormRow">
        <CardElement options={CARD_OPTIONS} onChange={onChange}/>
    </div>
);
