import type {
    FAQPage,
    Organization,
    Person,
    SoftwareApplication,
    WithContext,
  } from "schema-dts";
  
  export const gregoryBaranski: WithContext<Person> = {
    "@context": "https://schema.org",
    "@type": "Person",
    birthDate: "2004-09-19",
    birthPlace: "Poland",
    email: "me@gbaranski.com",
    gender: "Male",
    givenName: "Grzegorz",
    height: "188 cm",
  };
  
  export const organisationSchema: WithContext<Organization> = {
    "@context": "https://schema.org",
    "@type": "Organization",
    awards: "E(x)plory 2023, Young Innovator 2022/2023",
    contactPoint: {
      "@type": "ContactPoint",
      availableLanguage: {
        "@type": "Language",
        name: "English",
        alternateName: "en",
      },
      email: "contact@wheretopark.app",
    },
    email: "contact@wheretopark.app",
    founder: gregoryBaranski,
    foundingDate: "2022-05-01",
    logo: "https://wheretopark.app/favicon.ico",
  
    name: "Where To Park",
    url: "https://wheretopark.app",
  };
  
  export const softwareApplication: WithContext<SoftwareApplication> = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    applicationCategory: "DriverApplication",
    applicationSubCategory: "UtilitiesApplication",
    installUrl: "https://apps.apple.com/us/app/where-to-park/id6444453582",
    operatingSystem: "iOS 13",
    screenshot: "https://wheretopark.app/preview.webp",
    abstract: "Where To Park app helps you find parking spots near you.",
    audience: {
      "@type": "Audience",
      audienceType: "car owners",
      geographicArea: {
        "@type": "Country",
        name: "Poland, Scotland",
      },
    },
    keywords: "parking lot",
    publisher: organisationSchema,
    aggregateRating: {
      "@type": "AggregateRating",
      ratingValue: "5.0",
      ratingCount: "3",
    },
    offers: {
      "@type": "Offer",
      availability: "InStock",
      price: "0.00",
    },
  
    name: "Where To Park",
    url: "https://wheretopark.app",
  };
  
  export const frequentlyAskedQuestions: WithContext<FAQPage> = {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    mainEntity: [
      {
        "@type": "Question",
        name: "Is the app free?",
        acceptedAnswer: {
          "@type": "Answer",
          text:
            "<b>Yes</b>, the app is free to use and ad-free. You can download it from the App Store.",
        },
      },
      {
        "@type": "Question",
        name: "Is the app available for iOS?",
        acceptedAnswer: {
          "@type": "Answer",
          text:
            '<b>Yes</b>, you can download it from the <a href="https://apps.apple.com/us/app/where-to-park/id6444453582">App Store<a/> since 21st of November 2022',
        },
      },
      {
        "@type": "Question",
        name: "Is the app available for Android?",
        acceptedAnswer: {
          "@type": "Answer",
          text: "Not yet, but we're working on it!",
        },
      },
      {
        "@type": "Question",
        name: "What parking lots do you support?",
        acceptedAnswer: {
          "@type": "Answer",
          text:
            "Our AI-based system supports <b>all kinds</b> of parking lots.",
        },
      },
      {
        "@type": "Question",
        name: "What cities do you support?",
        acceptedAnswer: {
          "@type": "Answer",
          text:
            "As of September 2023 - we support Gdańsk, Gdynia, Sopot, Warszawa(Warsaw), Poznań, Kłodzko, Glasgow.",
        },
      },
      {
        "@type": "Question",
        name: "How to add a new parking lot to your app?",
        acceptedAnswer: {
          "@type": "Answer",
          text:
            "Contact us at contact@wheretopark.app with some basic details about your parking lot",
        },
      },
    ],
  };