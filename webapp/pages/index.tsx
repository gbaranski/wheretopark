import {parseParkingLots} from '../lib/types'
import {storekeeperClient} from '../lib/client'
import {GetStaticPropsContext, GetStaticPropsResult} from "next";
import {Home} from "../components/Home";

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


export async function getStaticProps(context: GetStaticPropsContext): Promise<GetStaticPropsResult<IndexProps>> {
    console.log("about to fetch parking lots")
    const parkingLots = await storekeeperClient.parkingLots()
    console.log({parkingLots})
    const props: IndexProps = {
        parkingLots: JSON.parse(parkingLots),
    }
    return {
        props, // will be passed to the page component as props
        revalidate: 10
    }
}

export default Index
