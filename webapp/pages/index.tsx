import {parseParkingLots} from '../lib/types'
import {storekeeperClient} from '../lib/client'
import {GetServerSidePropsContext} from "next";
import {Home} from "./Home";

type IndexProps = {
    parkingLots: any,
}

const Index = ({parkingLots: parkingLotsJSON}: IndexProps) => {
    const parkingLots = parseParkingLots(JSON.stringify(parkingLotsJSON))
    return (
        <>
            <Home parkingLots={parkingLots}/>
        </>
    )
}


export async function getServerSideProps(context: GetServerSidePropsContext) {
    const parkingLots = await storekeeperClient.parkingLots()
    context.res.setHeader(
        'Cache-Control',
        'public, s-maxage=10, stale-while-revalidate=59'
    )
    const props: IndexProps = {
        parkingLots: JSON.parse(parkingLots),
    }
    return {
        props, // will be passed to the page component as props
    }
}

export default Index
