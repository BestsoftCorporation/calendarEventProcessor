# CalendarEventsProcessor Service

This is a service which can be used to manage and process google calendar events. When endpoint /fetchEvents is 
triggered by GCS, CalendarHandler should fetch the latest events (by default from system1 prefix (ex 
system1@yourdomain.com)) and process them , by reading the content and update/add events in calendar and mongoDB.

This is an example of trip events processing in google calendar. User can define what type of transportion is used 
to commute.

#Stack
-Go 

-MongoDB

# Changelog
- **v1**:
- Calendar integration
- Google Maps integration
- Catching models

##RUN

Set credentials.json from GCS in **calendar/usecase/credentials.json**