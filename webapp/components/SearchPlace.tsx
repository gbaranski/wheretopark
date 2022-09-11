import * as React from 'react';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import Typography from '@mui/material/Typography';
import throttle from 'lodash/throttle';
import { MAPBOX_ACCESS_TOKEN } from '../environment';
import { Box, Button, CircularProgress, Grid } from '@mui/material';
import mbxGeocoding, { GeocodeFeature, GeocodeRequest, GeocodeResponse } from '@mapbox/mapbox-sdk/services/geocoding';
import { LocationOn, MyLocation } from '@mui/icons-material';
import { Coordinate } from '../lib/types'

type Props = {
  onSelect: (option: Coordinate | null) => void;
  buttonNeighbour: React.FC
};

const SearchPlace = ({ onSelect, buttonNeighbour }: Props) => {
  const [selected, setSelected] = React.useState<GeocodeFeature | null>(null);
  const [inputValue, setInputValue] = React.useState('');
  const [options, setOptions] = React.useState<readonly GeocodeFeature[]>([]);
  const [loading, setLoading] = React.useState(false);

  const fetch = React.useMemo(
    () =>
      throttle(
        (
          request: GeocodeRequest,
          callback: (response: GeocodeResponse) => void,
        ) => {
          const geocodingClient = mbxGeocoding({
            accessToken: MAPBOX_ACCESS_TOKEN,
          });
          geocodingClient.forwardGeocode(request).send().then((response) => {
            callback(response.body)
          });
        },
        1000,
      ),
    [],
  );

  React.useEffect(() => {
    setLoading(true)
    fetch({ query: inputValue, countries: ['pl'], limit: 5, autocomplete: true }, (response) => {
      response.features.sort((a, b) => a.relevance - b.relevance);
      console.log({ response });
      setOptions(response.features)
      setLoading(false)
    });
  }, [inputValue, fetch])

  const selectCurrentPosition = () => {
    navigator.geolocation.getCurrentPosition((position) => {
      onSelect(new Coordinate(position.coords.latitude, position.coords.longitude));
    }, (error) => {
      alert("Could not retrieve current position");
      console.error("getCurrentPosition error", error)
    });

  }

  return (
    <div style={{ marginInline: 10 }}>
      <Autocomplete
        id="mapbox-autocomplete"
        getOptionLabel={(option) => option.text}
        filterOptions={(x) => x}
        options={options}
        autoComplete
        includeInputInList
        filterSelectedOptions
        value={selected}
        isOptionEqualToValue={(option, value) => option.id === value.id}
        noOptionsText='No results'
        onChange={(event: any, newValue: GeocodeFeature | null) => {
          setOptions(newValue ? [newValue, ...options] : options);
          setSelected(newValue);
          onSelect(newValue ? new Coordinate(newValue.center[1], newValue.center[0]) : null);
        }}
        onInputChange={(event, newInputValue) => {
          setInputValue(newInputValue);
        }}
        renderInput={(params) => (
          <TextField
            {...params}
            value=""
            label="Add a location"
            fullWidth
            InputProps={{
              ...params.InputProps,
              type: 'search',
              endAdornment: (
                <React.Fragment>
                  {loading ? <CircularProgress color="inherit" size={20} /> : null}
                  {params.InputProps.endAdornment}
                </React.Fragment>
              ),
            }}
          />
        )}
        renderOption={(props, option) => {
          return (
            <li {...props} key={option.id}>
              <Grid container alignItems="center">
                <Grid item>
                  <Box
                    component={LocationOn}
                    sx={{ color: "text.secondary", mr: 2 }}
                  />
                </Grid>
                <Grid item xs>
                  <Typography variant="body1" color="text.primary">
                    {option.text}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {[option.properties.address, ...option.context.map((v) => v.text)].filter((s) => s).join(', ')}
                  </Typography>
                </Grid>
              </Grid>
            </li>
          );
        }}
      />
      <Box sx={{ 
        display: "inline-flex",
        flexDirection: "row",
        justifyContent: "space-between"
        
      }}>
        <Button color="secondary" startIcon={<MyLocation />} onClick={selectCurrentPosition}>Current location</Button>
        {buttonNeighbour({})}
      </Box>
    </div>
  );
}

export default SearchPlace;