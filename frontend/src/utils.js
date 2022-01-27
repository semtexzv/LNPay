import {decode} from "@node-lightning/invoice";

export function mapNetwork(net) {
    switch (net.toLowerCase()) {
        case "bc":
            return "Bitcoin Mainnet";
        case "tb":
            return "Bitcoin Testnet";
        case "sb":
            return "Bitcoin Simnet";
        default:
            return "unknown";
    }
}

export function decodeInvoice(invoice) {
    return function () {
        try {
            return decode(invoice);
        } catch (e) {
            return null;
        }
    }();
}

export function removePrefix(str, prefix) {
    const hasPrefix = str.indexOf(prefix) === 0;
    return hasPrefix ? str.substr(prefix.length) : str.toString();
}