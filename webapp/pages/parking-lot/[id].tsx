import {GetServerSidePropsContext, GetServerSidePropsResult} from "next";
import {storekeeperClient} from "../../lib/client";
import {ParkingLotID, parseParkingLots} from "../../lib/types";
import {Home} from "../Home";

type Props = {
    id: ParkingLotID
    parkingLots: any,
}

export default ({parkingLots: parkingLotsJSON, id}: Props) => {
    const parkingLots = parseParkingLots(JSON.stringify(parkingLotsJSON));
    return (
        <Home parkingLots={parkingLots} selectedParkingLotID={id}/>
    )

}

export async function getServerSideProps(context: GetServerSidePropsContext): Promise<GetServerSidePropsResult<Props>> {
    const {id} = context.query as { id: string };
    const parkingLots = await storekeeperClient.parkingLots();
    context.res.setHeader(
        'Cache-Control',
        'public, s-maxage=10, stale-while-revalidate=59'
    )
    const props: Props = {
        id,
        parkingLots: JSON.parse(parkingLots),
    }
    return {
        props, // will be passed to the page component as props
    }
}
