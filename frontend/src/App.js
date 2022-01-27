import './App.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faLinkedin, faTwitter} from '@fortawesome/free-brands-svg-icons';
import React, {useState, useEffect, useMemo} from 'react';

import {faAt} from '@fortawesome/free-solid-svg-icons';
import {faBe} from '@fortawesome/free-regular-svg-icons';
import {Elements} from '@stripe/react-stripe-js';
import {loadStripe} from '@stripe/stripe-js';
import {Flow} from "./flow/Flow";
import {API_URL} from "./config";

const stripePromise = loadStripe('INSERT HERE');

function NodeConnectLink() {
    const [metadata, setMetadata] = useState({});

    useEffect(() => {
        fetch(`${API_URL}/api/v1/metadata`)
            .then(r => r.json())
            .then(r => setMetadata(r)).catch(console.log);
    }, []);

    const idText = useMemo(() => `${metadata['pub_key'] ? metadata['pub_key'] : 'loading'}${metadata.ip ? "@" + metadata.ip : ""}`, [metadata]);
    const idLink = useMemo(() => {
        let base;
        if (metadata.network != 'mainnet') {
            base = `https://1ml.com/${metadata.network}/`;
        } else {
            base = `https://1ml.com/`;
        }
        return base + `node/${metadata.pub_key}`;
    });

    return (
        <>
            <div className="MetaSection">
                <div style={{textAlign: 'right', fontWeight: 'bold'}}>{metadata.network}</div>
                <a href={idLink} style={{textAlign: 'right'}}>{idText}</a>
            </div>
        </>
    );
}

function AuthorLinks() {
    return (<div className="SocialSection">
        <a className="SocialIcon" href="mailto:semtexzv@gmail.com"><FontAwesomeIcon icon={faAt}/></a>
        <a className="SocialIcon" href="https://twitter.com/Semtexzv"><FontAwesomeIcon icon={faTwitter}/></a>
        <a className="SocialIcon" href="https://www.linkedin.com/in/michalhornicky/"><FontAwesomeIcon
            icon={faLinkedin}/></a>
    </div>);
}

function App() {

    return (
        <>
            <div className={"AppWrapper"}>
                <div className="App">
                    <Elements stripe={stripePromise}>
                        <Flow/>
                    </Elements>
                </div>
            </div>
            <div className={"BottomSection"} style={{flex: '0'}}>
                <AuthorLinks/>
                <NodeConnectLink/>
            </div>
        </>
    );
}

export default App;
