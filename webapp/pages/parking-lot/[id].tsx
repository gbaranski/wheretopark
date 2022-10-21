import {GetServerSidePropsContext, GetServerSidePropsResult, GetStaticPathsContext, GetStaticPathsResult, GetStaticPropsContext, GetStaticPropsResult} from "next";
import {storekeeperClient} from "../../lib/client";
import {encodeParkingLots, ParkingLot, ParkingLotID, parseParkingLots, toRecord} from "../../lib/types";
import {Home} from "../../components/Home";

type Props = {
    id: ParkingLotID
    parkingLots: any,
}

export const ParkingLotByID = ({parkingLots: parkingLotsJSON, id}: Props) => {
    console.log(`front | encoded parking lots: ${JSON.stringify(parkingLotsJSON)}`);
    const parkingLots = parseParkingLots(JSON.stringify(parkingLotsJSON));
    return (
        <Home parkingLots={parkingLots} selectedParkingLotID={id}/>
    )
}

export async function getStaticProps(context: GetStaticPropsContext): Promise<GetStaticPropsResult<Props>> {
    const {id} = context.params as { id: string };
    console.log(`generating page for ${id}`);
    const parkingLots = encodeParkingLots(await storekeeperClient.parkingLots());
    console.log(`backend | encoded parking lots: ${parkingLots}`);
    const props: Props = {
        id,
        parkingLots: JSON.parse(parkingLots),
    }
    return {
        props, // will be passed to the page component as props
        // revalidate: 10,
    }
}

export async function getStaticPaths(context: GetStaticPathsContext): Promise<GetStaticPathsResult> {
    const parkingLotsMap = await storekeeperClient.parkingLots();
    const parkingLots = toRecord<ParkingLotID, ParkingLot>(parkingLotsMap);
    return {
      paths: Object.entries(parkingLots).map(([id, _]) => ({
        params: { id: id }
      })),
      fallback: true
    }
  }

export default ParkingLotByID