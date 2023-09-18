declare module "@recogito/annotorious" {
    declare function init(any);
  
    // {
    //   "@context": "http://www.w3.org/ns/anno.jsonld",
    //   "type": "Selection",
    //   "body": [
    //     {
    //       "type": "TextualBody",
    //       "purpose": "tagging",
    //       "value": "ParkingSpot"
    //     }
    //   ],
    //   "target": {
    //     "source": "http://localhost:5173/api/capture?src=https%3A%2F%2Fcam4out.klemit.net%2Fhls%2Fcamn826.m3u8",
    //     "selector": {
    //       "type": "SvgSelector",
    //       "value": "<svg><polygon points=\"279.5,177.5 351,177.5 367,204 278.5,235 278,235\"></polygon></svg>"
    //     }
    //   }
    // }
    type WebAnnotation = {
      "@context": string;
      type: "Annotation";
      body: {
        type: "TextualBody";
        purpose?: "tagging";
        value: string;
        format?: "text/plain";
      }[];
      id: string;
      target: {
        source: string;
        selector: {
          type: "SvgSelector";
          value: string;
        };
      };
    };
  
    export { init, WebAnnotation };
  }
  