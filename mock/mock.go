/*
Package mock contains mocked http endpoints to facilitate testing of the ourwx
package.
*/

package mock

import (
	"net/http"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const (
	htmlForecast = `<!DOCTYPE html>
<html class="no-js">
    <head>
        <!-- Meta -->
        <meta name="viewport" content="width=device-width">
        <link rel="schema.DC" href="http://purl.org/dc/elements/1.1/" /><title>National Weather Service</title><meta name="DC.title" content="National Weather Service" /><meta name="DC.description" content="NOAA National Weather Service National Weather Service" /><meta name="DC.creator" content="US Department of Commerce, NOAA, National Weather Service" /><meta name="DC.date.created" scheme="ISO8601" content="" /><meta name="DC.language" scheme="DCTERMS.RFC1766" content="EN-US" /><meta name="DC.keywords" content="weather, National Weather Service" /><meta name="DC.publisher" content="NOAA's National Weather Service" /><meta name="DC.contributor" content="National Weather Service" /><meta name="DC.rights" content="http://www.weather.gov/disclaimer.php" /><meta name="rating" content="General" /><meta name="robots" content="index,follow" />

        <!-- Icons -->
        <link rel="shortcut icon" href="./images/favicon.ico" type="image/x-icon" />

        <!-- CSS -->
        <link rel="stylesheet" href="css/bootstrap-3.2.0.min.css">
        <link rel="stylesheet" href="css/bootstrap-theme-3.2.0.min.css">
        <link rel="stylesheet" href="css/font-awesome-4.3.0.min.css">
        <link rel="stylesheet" href="css/ol-3.7.0.css" type="text/css">
        <link rel="stylesheet" type="text/css" href="css/mapclick.css" />
        <!--[if lte IE 7]><link rel="stylesheet" type="text/css" href="css/bootstrap-ie7.css" /><![endif]-->
        <!--[if lte IE 9]><link rel="stylesheet" type="text/css" href="css/mapclick-ie.css" /><![endif]-->
        <link rel="stylesheet" type="text/css" href="css/print.css" />
        <link rel="stylesheet" type="text/css" href="css/search.css" />

        <!-- Javascript -->
        <script type="text/javascript" src="js/lib/modernizr-2.8.3.js"></script>
        <script type="text/javascript" src="js/lib/json3-3.3.2.min.js"></script>
        <script type="text/javascript" src="js/lib/jquery-1.11.3.min.js"></script>
        <script type="text/javascript" src="js/lib/jquery.hoverIntent-1.8.1.min.js"></script>
        <script type="text/javascript" src="js/lib/bootstrap-3.2.0.min.js"></script>
        <script type="text/javascript" src="js/lib/ol-3.7.0.js"></script>
        <!--[if lte IE 8]><script type="text/javascript" src="js/respond.min.js"></script><![endif]-->
        <script type="text/javascript" src="js/jquery.autocomplete.min.js"></script>
        <script type="text/javascript" src="js/cfisurvey/cfi.js"></script>
        <script type="text/javascript" src="js/forecast.esri.js"></script>
        <script type="text/javascript" src="js/forecast.search.js"></script>
        <script type="text/javascript" src="js/forecast.openlayers.js"></script>
        <script type="text/javascript" src="js/browserSniffer.js"></script>
        <script type="text/javascript" src="js/federated-analytics.js"></script>
        <script type="text/javascript">
            (function (i, s, o, g, r, a, m) {
                i['GoogleAnalyticsObject'] = r;
                i[r] = i[r] || function () {
                    (i[r].q = i[r].q || []).push(arguments)
                }, i[r].l = 1 * new Date();
                a = s.createElement(o),
                        m = s.getElementsByTagName(o)[0];
                a.async = 1;
                a.src = g;
                m.parentNode.insertBefore(a, m)
            })(window, document, 'script', '//www.google-analytics.com/analytics.js', 'ga');

            ga('create', 'UA-40768555-1', 'weather.gov');
            ga('set', 'anonymizeIp', true);
            ga('require', 'linkid');
            ga('send', 'pageview');
        </script>

    </head>
    <body>
        <!-- DIV FOR CFI GROUP SURVEY::ALSO REQUIRES SCRIPT TAG IN HEADER -->
        <div id='ZN_9tslIS2mj3KoPgp'><!--DO NOT REMOVE-CONTENTS PLACED HERE--></div>

        <style>
            /* STYLE FOR DUAL ICON PREVIEW */
            .applicationnotificationContainerBanner {
                padding: 0 3rem 2rem 3rem;
                background: #fff;
                color: #555;
                margin-bottom: -.5rem;
                text-align: center;
                font-size: 1.2rem;
            }
            .applicationnotificationContainerBannerTeaser {
                display: inline-block;
                cursor: pointer;
            }
            .applicationnotificationContainerBannerTeaserIcon {
                float: left;
                width: 6.2rem;
                top:-3px;
                position: relative;
            }
            .applicationnotificationContainerBannerTeaserTitle {
                padding: 6px 0 0 0;
                font-weight: bold;
                font-size: 1.6rem;
                line-height: 1.6rem;
                margin-left: 7.2rem;
                text-align: left;
            }
            .applicationnotificationContainerBannerTeaserSubtitle {
                margin-left: 7.2rem;
                font-size: 1.2rem;
                line-height: 1.2rem;
                text-align: left;
                padding-top:.5rem;
            }
            .applicationnotificationContainerBannerDetails {
                display: none;
                clear: both;
                text-align: left;
                padding: 1rem 0;
                max-width: 750px;
                margin: 0 auto;
            }
            .applicationnotificationContainerBanner a {
                text-decoration: underline;
                padding-top:.3rem;
                display: block;
            }
            main.container {
                margin-top: -10px !important;
            }
            @media screen and (max-width:480px){
                .applicationnotificationContainerBanner{
                    padding:0 .8em 1em .8em;
                }
            }
        </style>
        <div class="applicationnotificationContainerBanner">
            <div id="applicationnotificationContainerButton-showDetails" class="applicationnotificationContainerBannerTeaser">
                <img src="images/applicationnotification.png" class="applicationnotificationContainerBannerTeaserIcon">
                <div class="applicationnotificationContainerBannerTeaserTitle">Notice of NWS' New Version of Forecast</div>
                <div class="applicationnotificationContainerBannerTeaserSubtitle">
                    A new version of Forecast will launch March 7, 2017.<br>
                    <a href="https://forecast-v3.weather.gov/documentation?redirect=legacy">Click here to visit the new site for details.</a><br> 
                </div>
            </div>
        </div>


        <main class="container">
            		<header class="row clearfix" id="page-header">
			<a href="http://www.noaa.gov" id="header-noaa" class="pull-left"><img src="/css/images/header_noaa.png" alt="National Oceanic and Atmospheric Administration"/></a>
			<a href="http://www.weather.gov" id="header-nws" class="pull-left"><img src="/css/images/header_nws.png" alt="National Weather Service"/></a>
			<a href="http://www.commerce.gov" id="header-doc" class="pull-right"><img src="/css/images/header_doc.png" alt="United States Department of Commerce"/></a>
		</header>
		
                    <nav class="navbar navbar-default row" role="navigation">
            <div class="container-fluid">
                <div class="navbar-header">
                    <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#top-nav">
                        <span class="sr-only">Toggle navigation</span>
                        <span class="icon-bar"></span>
                        <span class="icon-bar"></span>
                        <span class="icon-bar"></span>
                    </button>
                </div>
                <div class="collapse navbar-collapse" id="top-nav">
                    <ul class="nav navbar-nav">
                        <li><a href="http://www.weather.gov">HOME</a></li>
                        <li class="dropdown"><a href="http://www.weather.gov/forecastmaps" class="dropdown-toggle" data-toggle="dropdown">FORECAST&nbsp;<span class="caret"></span></a><ul class="dropdown-menu" role="menu"><li><a href="http://www.weather.gov">Local</a></li><li><a href="http://digital.weather.gov">Graphical</a></li><li><a href="http://www.aviationweather.gov/">Aviation</a></li><li><a href="http://www.nws.noaa.gov/om/marine/home.htm">Marine</a></li><li><a href="http://water.weather.gov/ahps/">Rivers and Lakes</a></li><li><a href="http://www.nhc.noaa.gov/">Hurricanes</a></li><li><a href="http://www.spc.noaa.gov/">Severe Weather</a></li><li><a href="http://www.srh.noaa.gov/ridge2/fire/">Fire Weather</a></li><li><a href="http://aa.usno.navy.mil/data/docs/RS_OneDay.php">Sun/Moon</a></li><li><a href="http://www.cpc.ncep.noaa.gov/">Long Range Forecasts</a></li><li><a href="http://www.cpc.ncep.noaa.gov">Climate Prediction</a></li></ul>                            </li>
                            <li class="dropdown"><a href="http://www.nws.noaa.gov/climate" class="dropdown-toggle" data-toggle="dropdown">PAST WEATHER&nbsp;<span class="caret"></span></a><ul class="dropdown-menu" role="menu"><li><a href="http://www.nws.noaa.gov/climate/">Past Weather</a></li><li><a href="http://www.nws.noaa.gov/climate/">Heating/Cooling Days</a></li><li><a href="http://www.nws.noaa.gov/climate/">Monthly Temperatures</a></li><li><a href="http://www.nws.noaa.gov/climate/">Records</a></li><li><a href="http://aa.usno.navy.mil/">Astronomical Data</a></li></ul>                            </li>
                            <li class="dropdown"><a href="http://www.weather.gov/safety" class="dropdown-toggle" data-toggle="dropdown">SAFETY&nbsp;<span class="caret"></span></a><ul class="dropdown-menu" role="menu"><li><a href="http://www.floodsafety.noaa.gov">Floods</a></li><li><a href="http://www.nws.noaa.gov/om/Tsunami/index.html">Tsunami</a></li><li><a href="http://www.nws.noaa.gov/beachhazards/">Beach Hazards</a></li><li><a href="http://www.nws.noaa.gov/om/fire/">Wildfire</a></li><li><a href="http://weather.gov/cold">Cold</a></li><li><a href="http://weather.gov/tornado">Tornadoes</a></li><li><a href="http://www.nws.noaa.gov/om/fog/">Fog</a></li><li><a href="http://www.nws.noaa.gov/airquality/">Air Quality</a></li><li><a href="http://www.nws.noaa.gov/om/heat/index.shtml">Heat</a></li><li><a href="http://www.nws.noaa.gov/om/hurricane/index.shtml">Hurricanes</a></li><li><a href="http://www.lightningsafety.noaa.gov/">Lightning</a></li><li><a href="http://www.ripcurrents.noaa.gov/">Rip Currents</a></li><li><a href="http://www.nws.noaa.gov/os/marine/safeboating/">Safe Boating</a></li><li><a href="http://weather.gov/thunderstorm">Thunderstorms</a></li><li><a href="http://www.nws.noaa.gov/om/space">Space Weather</a></li><li><a href="http://www.nws.noaa.gov/os/uv/">Sun (Ultraviolet Radiation)</a></li><li><a href="http://www.weather.gov/safetycampaign">Safety Campaigns</a></li><li><a href="http://www.weather.gov/wind">Wind</a></li><li><a href="http://www.weather.gov/om/drought/">Drought</a></li><li><a href="http://weather.gov/winter">Winter Weather</a></li></ul>                            </li>
                            <li class="dropdown"><a href="http://www.weather.gov/informationcenter" class="dropdown-toggle" data-toggle="dropdown">INFORMATION&nbsp;<span class="caret"></span></a><ul class="dropdown-menu" role="menu"><li><a href="http://www.weather.gov/Owlie's">Owlie's Kids Page</a></li><li><a href="http://www.nws.noaa.gov/com/weatherreadynation/wea.html">Wireless Emergency Alerts</a></li><li><a href="http://www.nws.noaa.gov/com/weatherreadynation">Weather-Ready Nation</a></li><li><a href="http://www.nws.noaa.gov/om/brochures.shtml">Brochures</a></li><li><a href="http://www.nws.noaa.gov/om/coop/">Cooperative Observers</a></li><li><a href="http://www.weather.gov/briefing/">Daily Briefing</a></li><li><a href="http://www.nws.noaa.gov/om/hazstats.shtml">Damage/Fatality/Injury Statistics</a></li><li><a href="http://mag.ncep.noaa.gov/">Forecast Models</a></li><li><a href="http://www.nws.noaa.gov/gis">GIS Data Portal</a></li><li><a href="http://www.nws.noaa.gov/nwr/">NOAA Weather Radio</a></li><li><a href="http://weather.gov/publications">Publications</a></li><li><a href="http://www.nws.noaa.gov/skywarn/">SKYWARN Storm Spotters</a></li><li><a href="http://www.nws.noaa.gov/stormready/">StormReady</a></li><li><a href="http://www.tsunamiready.noaa.gov">TsunamiReady</a></li></ul>                            </li>
                            <li class="dropdown"><a href="http://www.weather.gov/owlie" class="dropdown-toggle" data-toggle="dropdown">EDUCATION&nbsp;<span class="caret"></span></a><ul class="dropdown-menu" role="menu"><li><a href="http://www.nws.noaa.gov/com/weatherreadynation/force.html">Be A Force of Nature</a></li><li><a href="http://www.weather.gov/owlie">NWS Education Home</a></li></ul>                            </li>
                            <li class="dropdown"><a href="http://www.weather.gov/contact-media/" class="dropdown-toggle" data-toggle="dropdown">NEWS&nbsp;<span class="caret"></span></a><ul class="dropdown-menu" role="menu"><li><a href="http://www.weather.gov/news">NWS News</a></li><li><a href="http://www.nws.noaa.gov/com/weatherreadynation/calendar.html">Events</a></li><li><a href="http://www.weather.gov/socialmedia">Social Media</a></li><li><a href="http://www.nws.noaa.gov/om/brochures.shtml">Pubs/Brochures/Booklets </a></li><li><a href="http://www.nws.noaa.gov/pa/nws_contacts.php">NWS Media Contacts</a></li></ul>                            </li>
                            <li class="dropdown"><a href="http://www.weather.gov/search" class="dropdown-toggle" data-toggle="dropdown">SEARCH&nbsp;<span class="caret"></span></a><ul class="dropdown-menu" role="menu">                                <li><!-- Begin search code -->
                                    <div id="site-search">
                                        <form method="get" action="http://search.usa.gov/search" style="margin-bottom: 0; margin-top: 0;">
                                            <input type="hidden" name="v:project" value="firstgov" /> 
                                            <label for="query">Search For</label> 
                                            <input type="text" name="query" id="query" size="12" /> 
                                            <input type="submit" value="Go" />
                                            <p>
                                                <input type="radio" name="affiliate" checked="checked" value="nws.noaa.gov" id="nws" /> 
                                                <label for="nws" class="search-scope">NWS</label> 
                                                <input type="radio" name="affiliate" value="noaa.gov" id="noaa" /> 
                                                <label for="noaa" class="search-scope">All NOAA</label>
                                            </p>
                                        </form>
                                    </div>
                                </li>
                                </ul>                            </li>
                            <li class="dropdown"><a href="http://www.weather.gov/about" class="dropdown-toggle" data-toggle="dropdown">ABOUT&nbsp;<span class="caret"></span></a><ul class="dropdown-menu" role="menu"><li><a href="http://www.weather.gov/about">About NWS</a></li><li><a href="http://www.weather.gov/organization">Organization</a></li><li><a href="http://www.nws.noaa.gov/sp">Strategic Plan</a></li><li><a href="https://sites.google.com/a/noaa.gov/nws-insider/">For NWS Employees</a></li><li><a href="http://www.weather.gov/international/">International</a></li><li><a href="http://www.weather.gov/organization">National Centers</a></li><li><a href="http://www.nws.noaa.gov/tg">Products and Services</a></li><li><a href="http://www.weather.gov/careers/">Careers</a></li><li><a href="http://www.weather.gov/contact">Contact Us</a></li><li><a href="http://www.nws.noaa.gov/glossary">Glossary</a></li></ul>                            </li>
                                                </ul>
                </div>
            </div>
        </nav>
        
	    <div class="contentArea">
			<!-- Start Forecastsearch -->
	<div class="" id="fcst-search">
	    <form name="getForecast" id="getForecast" class="form-inline" role="form" action="http://forecast.weather.gov/zipcity.php" method="get">
		<div id="getfcst-body">
		    <input name="inputstring" type="text" class="form-control" id="inputstring" placeholder="" />
		    <input name="btnSearch" id="btnSearch" class="btn btn-default" type="submit" value="Go" />
		    <div id="txtHelp"><a href="javascript:void(window.open('http://weather.gov/ForecastSearchHelp.html','locsearchhelp','status=0,toolbar=0,location=0,menubar=0,directories=0,resizable=1,scrollbars=1,height=500,width=530').focus());">View Location Examples</a></div>
		</div>
		<div id="txtError">
		    <div id="errorNoResults" style="display:none;">Sorry, the location you searched for was not found. Please try another search.</div>
		    <div id="errorMultipleResults" style="display:none">Multiple locations were found. Please select one of the following:</div>
		    <div id="errorChoices" style="display:none"></div>
		    <input id="btnCloseError" type="button" value="Close" style="display:none" />
		</div>
		<div id="getfcst-head">
		    <p>Your local forecast office is</p>
		    <h3 id="getfcst-headOffice"></h3>
		</div>
	    </form>
	</div>
	<!-- end Forecastsearch -->
        
		<link rel="stylesheet" type="text/css" href="/css/topnews.css">
<div id="news-items">
    <div id="topnews">
    <div class="icon"><img src="/images/news-important.jpg"></div>
    <div class="body">
        <h1 style="font-size: 11pt;">Severe Weather Continues Today Across the Eastern U.S.; Very Wet in Hawaii</h1>
        <p>
            Severe thunderstorms continue ahead of a cold front stretching from the Deep South to the Mid-Atlantic and Northeast. Damaging winds will be the greatest threat from these storms, although, isolated tornadoes will be possible. Storms may also produce large hail and flash flooding. In the Pacific, a strong trough is bringing very heavy rainfall and the potential for flash flooding to Hawaii.  
            <a href="http://www.spc.noaa.gov/" target="_blank">Read More &gt;</a>
        </p>
    </div>
</div>

</div>
		<script type="text/javascript">(function ($) { var topnews = $("#topnews"); topnews.hide(); $.get("siteNews.php", {a:"pqr"},function(response){ if (response !== "false") topnews.replaceWith($(response)); topnews.show(); }); })(jQuery);</script><!-- PageFormat-Land -->
<script language=javascript>document.title = $('<div/>').html('7-Day Forecast for Latitude 45.53&deg;N and Longitude 122.67&deg;W (Elev. 200 ft)').text();</script><img src="images/track_land_point.png" style="display:none;" />
<div id="quickLinks">
	<span class="lang-spanish"><a href="http://forecast.weather.gov/MapClick.php?lat=45.52344714800046&lon=-122.67620703599971&lg=sp">En Espa&ntilde;ol</a></span>
	<div class="addthis_toolbox addthis_default_style addthis-forecast">
	    <a href="http://www.addthis.com/bookmark.php?v=250&amp;pubid=ra-5127a6364d551d04" class="addthis_button_compact">Share</a>
	    <span class="addthis_separator">|</span>
	    <a class="addthis_button_preferred_1"></a>
	    <a class="addthis_button_preferred_2"></a>
	    <a class="addthis_button_preferred_3"></a>
	    <a class="addthis_button_preferred_4"></a>
	    <a class="addthis_button_preferred_5"></a>
	</div>
	<script type="text/javascript">
		var addthis_config = addthis_config || {data_track_addressbar:true, pubid: 'xa-4b05b2d91f18c9cc'};
	    $(document).ready(function(){
			jQuery.ajax({
				url: "//s7.addthis.com/js/300/addthis_widget.js#async=1",
				dataType: "script",
				cache: false
			});
	    });
	</script>
</div>

<!-- Current Conditions -->
<div id="current-conditions" class="panel panel-default">

	<!-- Current Conditions header row -->
    <div class="panel-heading">
		<div>
		    <b>Current conditions at</b>
		    <h2 class="panel-title">Portland, Portland International Airport (KPDX)</h2>
		    <span class="smallTxt"><b>Lat:&nbsp;</b>45.59578&deg;N<b>Lon:&nbsp;</b>122.60917&deg;W<b>Elev:&nbsp;</b>20ft.</span>
	    </div>
    </div>
    <div class="panel-body" id="current-conditions-body">
		<!-- Graphic and temperatures -->
		<div id="current_conditions-summary" class="pull-left" >
		    		    <img src="newimages/large/ra.png" alt="" class="pull-left" />
		    		    <p class="myforecast-current">Lt Rain</p>
		    <p class="myforecast-current-lrg">45&deg;F</p>
		    <p class="myforecast-current-sm">7&deg;C</p>
		</div>
		<div id="current_conditions_detail" class="pull-left">
		    <table>
            <tr>
            <td class="text-right"><b>Humidity</b></td>
            <td>81%</td>
            </tr>
            <tr>
            <td class="text-right"><b>Wind Speed</b></td>
            <td>SSW 14 MPH</td>
            </tr>
            <tr>
            <td class="text-right"><b>Barometer</b></td>
            <td>30.47 in</td>
            </tr>
            <tr>
            <td class="text-right"><b>Dewpoint</b></td>
            <td>39&deg;F (4&deg;C)</td>
            </tr>
            <tr>
            <td class="text-right"><b>Visibility</b></td>
            <td>10.00 mi</td>
            </tr>
            <tr><td class="text-right"><b>Wind Chill</b></td><td>39&deg;F (4&deg;C)</td></tr>            <tr>
            <td class="text-right"><b>Last update</b></td>
            <td>
                01 Mar 8:35 am PST             </td>
            </tr>
		    </table>
		</div>
		<div id="current_conditions_station">
		    <div class="current-conditions-extra">
                            <!-- Right hand section -->
            <p class="moreInfo"><b>More Information:</b></p><p><a id="localWFO" href="http://www.wrh.noaa.gov/pqr" title="Portland, OR"><span class="hideText">Local</span> Forecast Office</a><a id="moreWx" href="http://www.wrh.noaa.gov/total_forecast/other_obs.php?wfo=pqr&zone=ORZ006">More Local Wx</a><a href="http://www.wrh.noaa.gov/mesowest/getobext.php?wfo=pqr&sid=KPDX&num=72&raw=0">3 Day History</a><a id="mobileWxLink" href="http://mobile.weather.gov/index.php?lat=45.5234&lon=-122.6762&unit=0&lg=english">Mobile Weather</a><a id="wxGraph" href="MapClick.php?lat=45.5234&lon=-122.6762&unit=0&amp;lg=english&amp;FcstType=graphical">Hourly <span class="hideText">Weather </span>Forecast</a></p>		    </div>
		<!-- /current_conditions_station -->
	    </div>
	    <!-- /current-conditions-body -->
	</div>
<!-- /Current Conditions -->
</div>

<!-- 7-Day Forecast -->
<div id="seven-day-forecast" class="panel panel-default">
    <div class="panel-heading">
	<b>Extended Forecast for</b>
	<h2 class="panel-title">
	    	    Portland OR	</h2>
    </div>
    <div class="panel-body" id="seven-day-forecast-body">
			<div id="seven-day-forecast-container"><ul id="seven-day-forecast-list" class="list-unstyled"><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Today<br><br></p>
<p><img src="newimages/medium/shra60.png" alt="Today: Showers likely, mainly before 10am.  Cloudy, with a high near 51. South southwest wind around 10 mph.  Chance of precipitation is 60%." title="Today: Showers likely, mainly before 10am.  Cloudy, with a high near 51. South southwest wind around 10 mph.  Chance of precipitation is 60%." class="forecast-icon"></p><p class="short-desc">Showers<br>Likely</p><p class="temp temp-high">High: 51 &deg;F</p></div></li><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Tonight<br><br></p>
<p><img src="DualImage.php?i=nshra&j=nfg&ip=20" alt="Tonight: A 20 percent chance of showers before 10pm.  Patchy fog after 10pm.  Otherwise, cloudy, with a low around 40. South wind around 6 mph. " title="Tonight: A 20 percent chance of showers before 10pm.  Patchy fog after 10pm.  Otherwise, cloudy, with a low around 40. South wind around 6 mph. " class="forecast-icon"></p><p class="short-desc">Slight Chance<br>Showers and<br>Patchy Fog<br>then Patchy<br>Fog</p><p class="temp temp-low">Low: 40 &deg;F</p></div></li><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Thursday<br><br></p>
<p><img src="newimages/medium/ra50.png" alt="Thursday: A 50 percent chance of rain.  Patchy fog before 10am.  Otherwise, cloudy, with a high near 50. South wind 7 to 10 mph.  New precipitation amounts between a tenth and quarter of an inch possible. " title="Thursday: A 50 percent chance of rain.  Patchy fog before 10am.  Otherwise, cloudy, with a high near 50. South wind 7 to 10 mph.  New precipitation amounts between a tenth and quarter of an inch possible. " class="forecast-icon"></p><p class="short-desc">Chance Rain<br>and Patchy<br>Fog</p><p class="temp temp-high">High: 50 &deg;F</p></div></li><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Thursday<br>Night</p>
<p><img src="newimages/medium/nra60.png" alt="Thursday Night: Rain likely, mainly after 4am.  Cloudy, with a low around 43. South wind 8 to 10 mph.  Chance of precipitation is 60%. New precipitation amounts between a tenth and quarter of an inch possible. " title="Thursday Night: Rain likely, mainly after 4am.  Cloudy, with a low around 43. South wind 8 to 10 mph.  Chance of precipitation is 60%. New precipitation amounts between a tenth and quarter of an inch possible. " class="forecast-icon"></p><p class="short-desc">Rain Likely</p><p class="temp temp-low">Low: 43 &deg;F</p></div></li><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Friday<br><br></p>
<p><img src="newimages/medium/ra90.png" alt="Friday: Rain.  High near 51. South wind 8 to 10 mph.  Chance of precipitation is 90%. New precipitation amounts between a quarter and half of an inch possible. " title="Friday: Rain.  High near 51. South wind 8 to 10 mph.  Chance of precipitation is 90%. New precipitation amounts between a quarter and half of an inch possible. " class="forecast-icon"></p><p class="short-desc">Rain</p><p class="temp temp-high">High: 51 &deg;F</p></div></li><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Friday<br>Night</p>
<p><img src="newimages/medium/nra90.png" alt="Friday Night: Rain before 10pm, then showers after 10pm.  Low around 40. Chance of precipitation is 90%." title="Friday Night: Rain before 10pm, then showers after 10pm.  Low around 40. Chance of precipitation is 90%." class="forecast-icon"></p><p class="short-desc">Rain</p><p class="temp temp-low">Low: 40 &deg;F</p></div></li><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Saturday<br><br></p>
<p><img src="newimages/medium/ra80.png" alt="Saturday: Rain.  High near 46. Chance of precipitation is 80%." title="Saturday: Rain.  High near 46. Chance of precipitation is 80%." class="forecast-icon"></p><p class="short-desc">Rain</p><p class="temp temp-high">High: 46 &deg;F</p></div></li><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Saturday<br>Night</p>
<p><img src="newimages/medium/nshra.png" alt="Saturday Night: Showers.  Cloudy, with a low around 37." title="Saturday Night: Showers.  Cloudy, with a low around 37." class="forecast-icon"></p><p class="short-desc">Showers</p><p class="temp temp-low">Low: 37 &deg;F</p></div></li><li class="forecast-tombstone">
<div class="tombstone-container">
<p class="period-name">Sunday<br><br></p>
<p><img src="newimages/medium/shra.png" alt="Sunday: Showers.  Cloudy, with a high near 44." title="Sunday: Showers.  Cloudy, with a high near 44." class="forecast-icon"></p><p class="short-desc">Showers</p><p class="temp temp-high">High: 44 &deg;F</p></div></li></ul></div>
<script type="text/javascript">
// equalize forecast heights
$(function () {
	var maxh = 0;
	$(".forecast-tombstone .short-desc").each(function () {
		var h = $(this).height();
		if (h > maxh) { maxh = h; }
	});
	$(".forecast-tombstone .short-desc").height(maxh);
});
</script>	</div>
</div>

<!-- Everything between 7-Day Forecast and Footer goes in this row -->
<div id="floatingDivs" class="row">
    <!-- Everything on the left-hand side -->
    <div class="col-md-7 col-lg-8">
        <!-- Detailed Forecast -->
        <div id="detailed-forecast" class="panel panel-default">
	    <div class="panel-heading">
            <h2 class="panel-title">Detailed Forecast</h2>
        </div>
	    <div class="panel-body" id="detailed-forecast-body">
            <div class="row row-odd row-forecast"><div class="col-sm-2 forecast-label"><b>Today</b></div><div class="col-sm-10 forecast-text">Showers likely, mainly before 10am.  Cloudy, with a high near 51. South southwest wind around 10 mph.  Chance of precipitation is 60%.</div></div><div class="row row-even row-forecast"><div class="col-sm-2 forecast-label"><b>Tonight</b></div><div class="col-sm-10 forecast-text">A 20 percent chance of showers before 10pm.  Patchy fog after 10pm.  Otherwise, cloudy, with a low around 40. South wind around 6 mph. </div></div><div class="row row-odd row-forecast"><div class="col-sm-2 forecast-label"><b>Thursday</b></div><div class="col-sm-10 forecast-text">A 50 percent chance of rain.  Patchy fog before 10am.  Otherwise, cloudy, with a high near 50. South wind 7 to 10 mph.  New precipitation amounts between a tenth and quarter of an inch possible. </div></div><div class="row row-even row-forecast"><div class="col-sm-2 forecast-label"><b>Thursday Night</b></div><div class="col-sm-10 forecast-text">Rain likely, mainly after 4am.  Cloudy, with a low around 43. South wind 8 to 10 mph.  Chance of precipitation is 60%. New precipitation amounts between a tenth and quarter of an inch possible. </div></div><div class="row row-odd row-forecast"><div class="col-sm-2 forecast-label"><b>Friday</b></div><div class="col-sm-10 forecast-text">Rain.  High near 51. South wind 8 to 10 mph.  Chance of precipitation is 90%. New precipitation amounts between a quarter and half of an inch possible. </div></div><div class="row row-even row-forecast"><div class="col-sm-2 forecast-label"><b>Friday Night</b></div><div class="col-sm-10 forecast-text">Rain before 10pm, then showers after 10pm.  Low around 40. Chance of precipitation is 90%.</div></div><div class="row row-odd row-forecast"><div class="col-sm-2 forecast-label"><b>Saturday</b></div><div class="col-sm-10 forecast-text">Rain.  High near 46. Chance of precipitation is 80%.</div></div><div class="row row-even row-forecast"><div class="col-sm-2 forecast-label"><b>Saturday Night</b></div><div class="col-sm-10 forecast-text">Showers.  Cloudy, with a low around 37.</div></div><div class="row row-odd row-forecast"><div class="col-sm-2 forecast-label"><b>Sunday</b></div><div class="col-sm-10 forecast-text">Showers.  Cloudy, with a high near 44.</div></div><div class="row row-even row-forecast"><div class="col-sm-2 forecast-label"><b>Sunday Night</b></div><div class="col-sm-10 forecast-text">Showers likely.  Cloudy, with a low around 38.</div></div><div class="row row-odd row-forecast"><div class="col-sm-2 forecast-label"><b>Monday</b></div><div class="col-sm-10 forecast-text">Showers.  Cloudy, with a high near 46.</div></div><div class="row row-even row-forecast"><div class="col-sm-2 forecast-label"><b>Monday Night</b></div><div class="col-sm-10 forecast-text">Showers.  Cloudy, with a low around 41.</div></div><div class="row row-odd row-forecast"><div class="col-sm-2 forecast-label"><b>Tuesday</b></div><div class="col-sm-10 forecast-text">Rain likely.  Cloudy, with a high near 50.</div></div>        </div>
	</div>
	<!-- /Detailed Forecast -->

        
        <!-- Additional Forecasts and Information -->
        <div id="additional_forecasts" class="panel panel-default">
	    <div class="panel-heading">
		<h2 class="panel-title">Additional Forecasts and Information</h2>
	    </div>

	    <div class="panel-body" id="additional-forecasts-body">
		<p class="myforecast-location"><a href="MapClick.php?zoneid=ORZ006">Zone Area Forecast for Greater Portland Metro Area, OR</a></p>
                <!-- First nine-ten links -->
		<div id="linkBlockContainer">
		    <div class="linkBlock">
                <ul class="list-unstyled">
                    <li><a href="http://forecast.weather.gov/product.php?site=PQR&issuedby=PQR&product=AFD&format=CI&version=1&glossary=1">Forecast Discussion</a></li>
                    <li><a href="MapClick.php?lat=45.5234&lon=-122.6762&unit=0&lg=english&FcstType=text&TextType=2">Printable Forecast</a></li>
                    <li><a href="MapClick.php?lat=45.5234&lon=-122.6762&unit=0&lg=english&FcstType=text&TextType=1">Text Only Forecast</a></li>
                </ul>
            </div>
		    <div class="linkBlock">
                <ul class="list-unstyled">
                    <li><a href="MapClick.php?lat=45.5234&lon=-122.6762&unit=0&lg=english&FcstType=graphical">Hourly Weather Forecast</a></li>
                    <li><a href="MapClick.php?lat=45.5234&lon=-122.6762&unit=0&lg=english&FcstType=digital">Tabular Forecast</a></li>
                    <!-- <li><a href="afm/PointClick.php?lat=45.5234&lon=-122.6762">Quick Forecast</a></li> -->
                </ul>
            </div>
		    <div class="linkBlock">
                <ul class="list-unstyled">
                    <li><a href="http://weather.gov/aq/probe_aq_data.php?latitude=45.5234&longitude=-122.6762">Air Quality Forecasts</a></li>
                    <li><a href="MapClick.php?lat=45.5234&lon=-122.6762&FcstType=text&unit=1&lg=en">International System of Units</a></li>
                    <li><a href="http://www.srh.weather.gov/srh/jetstream/webweather/pinpoint_max.htm">About Point Forecasts</a></li>
                                        <li><a href="http://www.wrh.noaa.gov/forecast/wxtables/index.php?lat=45.5234&lon=-122.6762">Forecast Weather Table Interface</a></li>
                                    </ul>
		    </div>
		    <!-- /First nine-ten links -->
                <!-- Additional links -->
                    <div class="linkBlock"><ul class="list-unstyled"><li><a href="http://www.wrh.noaa.gov/total_forecast/getprod.php?wfo=pqr&pil=PFM&sid=PQR" target="_self">PFM (Forecast Matrix)</a></li><li><a href="http://www.wrh.noaa.gov/pqr/info/pdf/pfm.pdf" target="_self">PFM Decoding Guide</a></li><li><a href="http://www.wrh.noaa.gov/pqr/rain.php" target="_self">Rainfall Forecasts</a></li><li><a href="http://www.wrh.noaa.gov/total_forecast/getprod.php?wfo=pqr&pil=RVS&sid=PQR" target="_self">Current River Levels</a></li></ul></div><div class="linkBlock"><ul class="list-unstyled"><li><a href="http://forecast.weather.gov/wxplanner.php?site=pqr" target="_self">Weather Planner</a></li><li><a href="http://www.wrh.noaa.gov/mesowest/frame.php?map=pqr" target="_self">Mapped Observations</a></li><li><a href="http://www.wrh.noaa.gov/total_forecast/getprod.php?wfo=pqr&pil=STO&sid=OR" target="_self">Road Conditions</a></li><li><a href="http://www.wrh.noaa.gov/total_forecast/getprod.php?wfo=pqr&pil=SAB&sid=SEA" target="_self">Avalanche Outlooks</a></li></ul></div><div class="linkBlock"><ul class="list-unstyled"><li><a href="http://www.wrh.noaa.gov/pqr/marine.php" target="_self">Marine Weather</a></li><li><a href="http://www.wrh.noaa.gov/firewx/index.php?wfo=pqr" target="_self">Fire Weather</a></li><li><a href="http://www.nws.noaa.gov/wtf/udaf/area/?site=pqr" target="_self">User Defined Area</a></li></ul></div>
		</div> <!-- /linkBlockContainer -->
	    </div><!-- /additional-forecasts-body-->
	</div> <!-- /additional_forecasts -->
    </div> <!-- /Everything on the left-hand side -->

    <!-- right-side-data -->
    <div class="col-md-5 col-lg-4" id="right-side-data">
	<div id="mapAndDescriptionArea">
        <!-- openlayer map -->
            <style>
#custom-search{
display: block;
position: relative;
z-index: 50;
top: 52px;
left: 60px;
}
#esri-geocoder-search{
display: block;
position: relative;
z-index: 50;
top: 52px;
left: 60px;
}
#emap{
margin-top:15px;
cursor:pointer;
height:370px;
width:100%;
border: 1px solid #ccc;
border-radius: 3px;
}
#switch-basemap-container{
}
#basemap-selection-form ul{
list-style: none;
 margin: 0px;
}
#basemap-selection-form li{
float: left;
}
.disclaimer{
margin-top:350px;
margin-left: 5px;
z-index: 100;
position: absolute;
text-transform: none;
}
.esriAttributionLastItem{
text-transform: none;
}
.esriSimpleSlider div{
height:22px;
line-height:20px;
width:20px;
}
#point-forecast-map-label {
text-align:center;
font-weight:bold;
color:black;
}
@media (max-width: 767px) {
#emap{
margin-top:.5em;
height:270px;
}
.disclaimer{
margin-top:250px;
}
}
</style>
<!-- forecast-map -->
<div class='point-forecast-map'>
    <div class='point-forecast-map-header text-center'>
        <div id="toolbar">
    	<div id="switch-basemap-container">
    	    <div id="basemap-selection-form" title="Choose a Basemap">
    		<div id="basemap-menu">
    		    <select name="basemap-selected" id="basemap-selected" autocomplete="off" title="Basemap Dropdown Menu">
    		    <option value="none">Select Basemap</option>
    		    <option value="topo" selected>Topographic</option>
    		    <option value="streets">Streets</option>
    		    <option value="satellite">Satellite</option>
    		    <option value="ocean">Ocean</option>
    		    </select>
    		</div>
    	    </div>
    	    <div id="point-forecast-map-label">
                    Click Map For Forecast
                </div>
    	</div><!-- //#switch-basemap-container -->
    	<div style="clear:both;"></div>
        </div><!-- //#toolbar -->
    </div><!-- //.point-forecast-map-header -->

    <div id="emap">
        <noscript><center><br><br><b>Map function requires Javascript and a compatible browser.</b></center></noscript>
        <div class="disclaimer"><a href='http://www.weather.gov/disclaimer#esri'>Disclaimer</a></div>
    </div><!-- //#emap -->

    <div class="point-forecast-map-footer">
        <img src="./images/wtf/maplegend_forecast-area.gif" width="100" height="16" alt="Map Legend">
    </div><!-- //.point-forecast-map-footer -->

</div> <!-- //.point-forecast-map -->
<!-- //forecast-map -->
        <!-- //openlayer map -->

	    <!-- About this Forecast -->
        <div id="about_forecast">
            <div class="fullRow">
                <div class="left">Point Forecast:</div>
                <div class="right">Portland OR<br>&nbsp;45.53&deg;N 122.67&deg;W (Elev. 200 ft)</div>
                    </div>
            <div class="fullRow">
                <div class="left"><a target="_blank" href="http://www.weather.gov/glossary/index.php?word=Last+update">Last Update</a>: </div>
                <div class="right">3:14 am PST Mar 1, 2017</div>
            </div>
            <div class="fullRow">
                <div class="left"><a target="_blank" href="http://www.weather.gov/glossary/index.php?word=forecast+valid+for">Forecast Valid</a>: </div>
                <div class="right">8am PST Mar 1, 2017-6pm PST Mar 7, 2017</div>
            </div>
            <div class="fullRow">
                <div class="left">&nbsp;</div>
                <div class="right"><a href="http://forecast.weather.gov/product.php?site=PQR&issuedby=PQR&product=AFD&format=CI&version=1&glossary=1">Forecast Discussion</a></div>
            </div>
            <div class="fullRow">
                <div class="left">&nbsp;</div>
                <div class="right">
                    <a href="MapClick.php?lat=45.5234&lon=-122.6762&unit=0&lg=english&FcstType=kml"><img src="/images/wtf/kml_badge.png" width="45" height="17" alt="Get as KML" /></a>
                    <a href="MapClick.php?lat=45.5234&lon=-122.6762&unit=0&lg=english&FcstType=dwml"><img src="/images/wtf/xml_badge.png" width="45" height="17" alt="Get as XML" /></a>
                </div>
            </div>
        </div>
	    <!-- /About this Forecast -->
	</div>
    
        <!--additionalForecast-->
        <div class="panel panel-default" id="additionalForecast">
            <div class="panel-heading">
                <h2 class="panel-title">Additional Resources</h2>
            </div>
            <div class="panel-body">

                <!-- Radar & Satellite Images -->
                <div id="radar" class="subItem">
                    <h4>Radar &amp; Satellite Image</h4>
                    <a href="http://radar.weather.gov/radar.php?rid=rtx&product=N0R&overlay=11101111&loop=no"><img src="http://radar.weather.gov/Thumbs/RTX_Thumb.gif" class="radar-thumb" alt="Link to Local Radar Data" title="Link to Local Radar Data"></a>                    <a href="http://www.wrh.noaa.gov/satellite/?wfo=pqr"><img src="http://www.ssd.noaa.gov/goes/west/wfo/pqr/ft.jpg" class="satellite-thumb" alt="Link to Satellite Data" title="Link to Satellite Data"></a>                </div>
                <!-- /Radar & Satellite Images -->
                <!-- Hourly Weather Forecast -->
                <div id="feature" class="subItem">
                    <h4>Hourly Weather Forecast</h4>
                    <a href="MapClick.php?lat=45.5234&lon=-122.6762&unit=0&lg=english&FcstType=graphical"><img src="newimages/medium/hourlyweather.png" class="img-responsive" /></a>
                </div>
                <!-- /Hourly Weather Forecast -->
                <!-- NDFD -->
                <div id="NDFD" class="subItem">
                    <h4>National Digital Forecast Database</h4>
                    <div class="one-sixth-first"><a href="http://graphical.weather.gov/sectors/pacnorthwest.php?element=MaxT"><img src="http://www.weather.gov/forecasts/graphical/images/thumbnail/latest_MaxMinT_pacnorthwest_thumbnail.png" border="0" alt="National Digital Forecast Database Maximum Temperature Forecast" title="National Digital Forecast Database Maximum Temperature Forecast" width="147" height="150"></a>
	 			<p><a href="http://graphical.weather.gov/sectors/pacnorthwest.php?element=MaxT">High Temperature</a></p></div><div class="one-sixth-first"><a href="http://graphical.weather.gov/sectors/pacnorthwest.php?element=Wx"><img src="http://www.weather.gov/forecasts/graphical/images/thumbnail/latest_Wx_pacnorthwest_thumbnail.png" border="0" alt="National Digital Forecast Database Weather Element Forecast" title="National Digital Forecast Database Weather Element Forecast" width="147" height="150"></a>
	 			<p><a href="http://graphical.weather.gov/sectors/pacnorthwest.php?element=Wx">Chance of Precipitation</a></p></div>                </div>
                <!-- /NDFD -->
            </div>
        </div>
        <!-- /additionalForecast -->

    </div>
    <!-- /col-md-4 -->
    <!-- /right-side-data -->
    <script language='javascript'>$( document ).ready(function() { load_openlayers_map('', '', '', '{"centroid_lat":"45.5234","centroid_lon":"-122.6762","lat1":"45.514","lon1":"-122.682","lat2":"45.5345","lon2":"-122.688","lat3":"45.539","lon3":"-122.658","lat4":"45.5185","lon4":"-122.652"}') });</script></div>
<!-- /row  -->


</div>
<!-- /PageFormat-Land -->

	    </div>
            <footer>
                        <div id="sitemap" class="sitemap-content row">
            <div class="col-xs-12">
                <div class="sitemap-columns">
                                                    <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href=" http://www.nws.noaa.gov/climate/">PAST WEATHER</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-1">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-1">
                                        <ul class="list-unstyled">
                                                                                            <li><a href=" http://www.nws.noaa.gov/climate/">Past Weather </a></li>
                                                                                                <li><a href=" http://www.cpc.ncep.noaa.gov/products/MD_index.shtml">Climate Monitoring </a></li>
                                                                                                <li><a href=" http://www.nws.noaa.gov/climate/">Heating/Cooling Days </a></li>
                                                                                                <li><a href=" http://www.nws.noaa.gov/climate/">Monthly Temps </a></li>
                                                                                                <li><a href=" http://www.nws.noaa.gov/climate/">Records </a></li>
                                                                                                <li><a href=" http://aa.usno.navy.mil/">Astronomical Data </a></li>
                                                                                                <li><a href="http://www.ncdc.noaa.gov/oa/mpp/">Certified Weather Data </a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                                <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href="http://alerts.weather.gov">ACTIVE ALERTS</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-2">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-2">
                                        <ul class="list-unstyled">
                                                                                            <li><a href=" http://alerts.weather.gov">Warnings By State</a></li>
                                                                                                <li><a href=" http://www.wpc.ncep.noaa.gov/ww.shtml">Excessive Rainfall and Winter Weather Forecasts</a></li>
                                                                                                <li><a href="http://water.weather.gov/ahps/?current_color=flood&amp;current_type=all&amp;fcst_type=obs&amp;conus_map=d_map">River Flooding </a></li>
                                                                                                <li><a href=" http://www.weather.gov">Latest Warnings</a></li>
                                                                                                <li><a href=" http://www.spc.noaa.gov/products/outlook/">Thunderstorm/Tornado Outlook </a></li>
                                                                                                <li><a href=" http://www.nhc.noaa.gov/">Hurricanes </a></li>
                                                                                                <li><a href=" http://www.spc.noaa.gov/products/fire_wx/">Fire Weather Outlooks </a></li>
                                                                                                <li><a href=" http://www.cpc.ncep.noaa.gov/products/stratosphere/uv_index/uv_alert.shtml">UV Alerts </a></li>
                                                                                                <li><a href=" http://www.drought.gov/">Drought </a></li>
                                                                                                <li><a href="http://www.swpc.noaa.gov/products/alerts-watches-and-warnings">Space Weather </a></li>
                                                                                                <li><a href=" http://www.nws.noaa.gov/nwr/">NOAA Weather Radio </a></li>
                                                                                                <li><a href=" http://alerts.weather.gov/">NWS CAP Feeds </a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                                <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href="http://www.weather.gov/current">CURRENT CONDITIONS</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-3">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-3">
                                        <ul class="list-unstyled">
                                                                                            <li><a href=" http://www.weather.gov/Radar">Radar </a></li>
                                                                                                <li><a href="http://www.cpc.ncep.noaa.gov/products/monitoring_and_data/">Climate Monitoring </a></li>
                                                                                                <li><a href=" http://water.weather.gov/ahps/">River Levels </a></li>
                                                                                                <li><a href=" http://water.weather.gov/precip/">Observed Precipitation </a></li>
                                                                                                <li><a href="http://www.nws.noaa.gov/om/osd/portal.shtml">Surface Weather </a></li>
                                                                                                <li><a href="ftp://tgftp.nws.noaa.gov/fax/barotrop.shtml">Upper Air </a></li>
                                                                                                <li><a href=" http://www.ndbc.noaa.gov/">Marine and Buoy Reports </a></li>
                                                                                                <li><a href="http://www.nohrsc.noaa.gov/interactive/html/map.html">Snow Cover </a></li>
                                                                                                <li><a href=" http://www.goes.noaa.gov">Satellite </a></li>
                                                                                                <li><a href=" http://www.swpc.noaa.gov/">Space Weather </a></li>
                                                                                                <li><a href="http://www.weather.gov/pr">International Observations</a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                                <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href="http://weather.gov/forecastmaps">FORECAST</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-4">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-4">
                                        <ul class="list-unstyled">
                                                                                            <li><a href=" http://www.weather.gov/">Local Forecast </a></li>
                                                                                                <li><a href="http://www.weather.gov/pr">International Forecasts</a></li>
                                                                                                <li><a href=" http://www.spc.noaa.gov/">Severe Weather </a></li>
                                                                                                <li><a href=" http://www.wpc.ncep.noaa.gov/">Current Outlook Maps </a></li>
                                                                                                <li><a href="http://www.cpc.ncep.noaa.gov/products/Drought">Drought </a></li>
                                                                                                <li><a href="http://www.weather.gov/fire">Fire Weather </a></li>
                                                                                                <li><a href=" http://www.wpc.ncep.noaa.gov/">Fronts/Precipitation Maps </a></li>
                                                                                                <li><a href=" http://www.nws.noaa.gov/forecasts/graphical/">Current Graphical Forecast Maps </a></li>
                                                                                                <li><a href="http://water.weather.gov/ahps/forecasts.php">Rivers </a></li>
                                                                                                <li><a href=" http://www.nws.noaa.gov/om/marine/home.htm">Marine </a></li>
                                                                                                <li><a href="http://www.opc.ncep.noaa.gov/marine_areas.php">Offshore and High Seas</a></li>
                                                                                                <li><a href=" http://www.nhc.noaa.gov/">Hurricanes </a></li>
                                                                                                <li><a href=" http://aviationweather.gov">Aviation Weather </a></li>
                                                                                                <li><a href="http://www.cpc.ncep.noaa.gov/products/OUTLOOKS_index.shtml">Climatic Outlook </a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                                <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href="http://www.weather.gov/informationcenter">INFORMATION CENTER</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-5">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-5">
                                        <ul class="list-unstyled">
                                                                                            <li><a href=" http://www.spaceweather.gov">Space Weather </a></li>
                                                                                                <li><a href="http://www.weather.gov/briefing/">Daily Briefing </a></li>
                                                                                                <li><a href=" http://www.nws.noaa.gov/om/marine/home.htm">Marine </a></li>
                                                                                                <li><a href="http://www.nws.noaa.gov/climate">Climate </a></li>
                                                                                                <li><a href="http://www.weather.gov/fire">Fire Weather </a></li>
                                                                                                <li><a href=" http://www.aviationweather.gov/">Aviation </a></li>
                                                                                                <li><a href="http://mag.ncep.noaa.gov/">Forecast Models </a></li>
                                                                                                <li><a href="http://water.weather.gov/ahps/">Water </a></li>
                                                                                                <li><a href="http://www.nws.noaa.gov/gis">GIS</a></li>
                                                                                                <li><a href="http://www.weather.gov/pr">International Weather</a></li>
                                                                                                <li><a href=" http://www.nws.noaa.gov/om/coop/">Cooperative Observers </a></li>
                                                                                                <li><a href="http://www.nws.noaa.gov/skywarn/">Storm Spotters </a></li>
                                                                                                <li><a href="http://www.tsunami.gov">Tsunami</a></li>
                                                                                                <li><a href="http://www.economics.noaa.gov">Facts and Figures </a></li>
                                                                                                <li><a href="http://water.noaa.gov/">National Water Center</a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                                <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href="http://weather.gov/safety">WEATHER SAFETY</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-6">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-6">
                                        <ul class="list-unstyled">
                                                                                            <li><a href="http://www.weather.gov/nwr/">NOAA Weather Radio</a></li>
                                                                                                <li><a href="http://www.weather.gov/stormready/">StormReady</a></li>
                                                                                                <li><a href="http://www.nws.noaa.gov/om/heat/index.shtml">Heat </a></li>
                                                                                                <li><a href=" http://www.lightningsafety.noaa.gov/">Lightning </a></li>
                                                                                                <li><a href=" http://www.nhc.noaa.gov/prepare/">Hurricanes </a></li>
                                                                                                <li><a href=" http://www.weather.gov/om/severeweather/index.shtml">Thunderstorms </a></li>
                                                                                                <li><a href=" http://www.weather.gov/om/severeweather/index.shtml">Tornadoes </a></li>
                                                                                                <li><a href=" http://www.weather.gov/om/severeweather/index.shtml">Severe Weather </a></li>
                                                                                                <li><a href=" http://www.ripcurrents.noaa.gov/">Rip Currents </a></li>
                                                                                                <li><a href="http://www.nws.noaa.gov/os/marine/safeboating/">Safe Boating</a></li>
                                                                                                <li><a href=" http://www.weather.gov/om/severeweather/index.shtml">Floods </a></li>
                                                                                                <li><a href=" http://www.weather.gov/om/winter/index.shtml">Winter Weather </a></li>
                                                                                                <li><a href=" http://www.weather.gov/os/uv/">Ultra Violet Radiation </a></li>
                                                                                                <li><a href=" http://www.weather.gov/airquality/">Air Quality </a></li>
                                                                                                <li><a href=" http://www.weather.gov/om/hazstats.shtml">Damage/Fatality/Injury Statistics </a></li>
                                                                                                <li><a href=" http://www.redcross.org/">Red Cross </a></li>
                                                                                                <li><a href=" http://www.fema.gov/">Federal Emergency Management Agency (FEMA) </a></li>
                                                                                                <li><a href=" http://www.weather.gov/om/brochures.shtml">Brochures </a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                                <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href="http://weather.gov/news">NEWS</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-7">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-7">
                                        <ul class="list-unstyled">
                                                                                            <li><a href=" http://weather.gov/news">Newsroom</a></li>
                                                                                                <li><a href=" http://weather.gov/socialmedia">Social Media </a></li>
                                                                                                <li><a href="http://www.nws.noaa.gov/com/weatherreadynation/calendar.html">Events</a></li>
                                                                                                <li><a href=" http://www.weather.gov/om/brochures.shtml">Pubs/Brochures/Booklets </a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                                <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href="http://weather.gov/owlie">EDUCATION</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-8">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-8">
                                        <ul class="list-unstyled">
                                                                                            <li><a href="http://weather.gov/owlie">NWS Education Home</a></li>
                                                                                                <li><a href="http://www.nws.noaa.gov/com/weatherreadynation/force.html">Be A Force of Nature</a></li>
                                                                                                <li><a href=" http://www.education.noaa.gov/Weather_and_Atmosphere/">NOAA Education Resources </a></li>
                                                                                                <li><a href=" http://www.weather.gov/glossary/">Glossary </a></li>
                                                                                                <li><a href=" http://www.srh.noaa.gov/srh/jetstream/">JetStream </a></li>
                                                                                                <li><a href=" http://www.weather.gov/training/">NWS Training Portal </a></li>
                                                                                                <li><a href=" http://www.lib.noaa.gov/">NOAA Library </a></li>
                                                                                                <li><a href="http://weather.gov/owlie">For Students, Parents and Teachers</a></li>
                                                                                                <li><a href="http://www.weather.gov/owlie/publication_brochures">Brochures </a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                                <div class="sitemap-section">
                                    <div class="panel-heading">
                                        <a class="sitemap-section-heading" href="http://weather.gov/about">ABOUT</a>
                                        <button type="button" class="menu-toggle pull-right" data-toggle="collapse" data-target="#sitemap-9">
                                            <span class="sr-only">Toggle menu</span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                            <span class="icon-bar"></span>
                                        </button>
                                    </div>
                                    <div class="sitemap-section-body panel-body collapsable collapse" id="sitemap-9">
                                        <ul class="list-unstyled">
                                                                                            <li><a href="http://weather.gov/organization">Organization </a></li>
                                                                                                <li><a href=" http://www.weather.gov/sp/">Strategic Plan </a></li>
                                                                                                <li><a href="https://sites.google.com/a/noaa.gov/nws-best-practices/">For NWS Employees </a></li>
                                                                                                <li><a href=" http://www.weather.gov/ia/home.htm">International </a></li>
                                                                                                <li><a href="http://www.ncep.noaa.gov/">National Centers </a></li>
                                                                                                <li><a href=" http://www.weather.gov/tg/">Products and Services </a></li>
                                                                                                <li><a href="http://www.weather.gov/careers/">Careers</a></li>
                                                                                                <li><a href=" http://www.weather.gov/glossary/">Glossary </a></li>
                                                                                                <li><a href="http://weather.gov/contact">Contact Us </a></li>
                                                                                        </ul>
                                    </div>
                                </div>
                                                </div>
            </div>
        </div>
        
                <!-- legal footer area -->
                		<div class="footer-legal">
			<div id="footerLogo" class="col-xs-12 col-sm-2 col-md-2">
				<a href="http://www.usa.gov"><img src="/css/images/usa_gov.png" alt="usa.gov" width="110" height="30" /></a>
			</div>
			<div class="col-xs-12 col-sm-4 col-md-4">
				<ul class="list-unstyled footer-legal-content">
				<li><a href="http://www.commerce.gov">US Dept of Commerce</a></li>
				<li><a href="http://www.noaa.gov">National Oceanic and Atmospheric Administration</a></li>
				<li><a href="http://www.weather.gov">National Weather Service</a></li>
				<li><a href="http://www.weather.gov/pqr">Portland, OR</a></li><li><br /><a href="mailto:w-pqr.webmaster@noaa.gov">Comments? Questions? Please Contact Us.</a></li>			</ul>
			</div>
			<div class="col-xs-12 col-sm-3 col-md-3">
				<ul class="list-unstyled">
					<li><a href="http://www.weather.gov/disclaimer">Disclaimer</a></li>
					<li><a href="http://www.cio.noaa.gov/services_programs/info_quality.html">Information Quality</a></li>
					<li><a href="http://www.weather.gov/help">Help</a></li>
					<li><a href="http://www.weather.gov/glossary">Glossary</a></li>
				</ul>
			</div>
			<div class="col-xs-12 col-sm-3 col-md-3">
				<ul class="list-unstyled">
					<li><a href="http://www.weather.gov/privacy">Privacy Policy</a></li>
					<li><a href="http://www.rdc.noaa.gov/~foia">Freedom of Information Act (FOIA)</a></li>
					<li><a href="http://www.weather.gov/about">About Us</a></li>
					<li><a href="http://www.weather.gov/careers">Career Opportunities</a></li>
				</ul>
			</div>
		</div>
		
            </footer>
        </main>
    </body>
</html>`

	observation = `<?xml version="1.0" encoding="ISO-8859-1"?> 
<?xml-stylesheet href="latest_ob.xsl" type="text/xsl"?>
<current_observation version="1.0"
	 xmlns:xsd="http://www.w3.org/2001/XMLSchema"
	 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	 xsi:noNamespaceSchemaLocation="http://www.weather.gov/view/current_observation.xsd">
	<credit>NOAA's National Weather Service</credit>
	<credit_URL>http://weather.gov/</credit_URL>
	<image>
		<url>http://weather.gov/images/xml_logo.gif</url>
		<title>NOAA's National Weather Service</title>
		<link>http://weather.gov</link>
	</image>
	<suggested_pickup>15 minutes after the hour</suggested_pickup>
	<suggested_pickup_period>60</suggested_pickup_period>
	<location>Portland, Portland International Airport, OR</location>
	<station_id>KPDX</station_id>
	<latitude>45.59578</latitude>
	<longitude>-122.60917</longitude>
	<elevation>0</elevation>
	<observation_time>Last Updated on Feb 27 2017, 8:53 am PST</observation_time>
    <observation_time_rfc822>Mon, 27 Feb 2017 08:53:00 -0800</observation_time_rfc822>
	<weather>Overcast</weather>
	<temperature_string>38.0 F (3.3 C)</temperature_string>
	<temp_f>38.0</temp_f>
	<temp_c>3.3</temp_c>
	<relative_humidity>86</relative_humidity>
	<wind_string>Southwest at 6.9 MPH (6 KT)</wind_string>
	<wind_dir>Southwest</wind_dir>
	<wind_degrees>230</wind_degrees>
	<wind_mph>6.9</wind_mph>
	<wind_kt>6</wind_kt>
	<wind_gust_mph>200</wind_gust_mph>
	<wind_gust_kt>173.8</wind_gust_kt>
	<pressure_string>1009.9 mb</pressure_string>
	<pressure_mb>1009.9</pressure_mb>
	<pressure_in>29.82</pressure_in>
	<dewpoint_string>34.0 F (1.1 C)</dewpoint_string>
	<dewpoint_f>34.0</dewpoint_f>
	<dewpoint_c>1.1</dewpoint_c>
	<windchill_string>33 F (1 C)</windchill_string>
    <windchill_f>33</windchill_f>
    <windchill_c>1</windchill_c>
    <heat_index_string>38.0 F (3.3 C)</heat_index_string>
    <heat_index_f>38.0</heat_index_f>
    <heat_index_c>3.3</heat_index_c>
	<visibility_mi>10.00</visibility_mi>
 	<icon_url_base>http://forecast.weather.gov/images/wtf/small/</icon_url_base>
	<two_day_history_url>http://www.weather.gov/data/obhistory/KPDX.html</two_day_history_url>
	<icon_url_name>ovc.png</icon_url_name>
	<ob_url>http://www.weather.gov/data/METAR/KPDX.1.txt</ob_url>
	<disclaimer_url>http://weather.gov/disclaimer.html</disclaimer_url>
	<copyright_url>http://weather.gov/disclaimer.html</copyright_url>
	<privacy_policy_url>http://weather.gov/notice.html</privacy_policy_url>
</current_observation>`

	alertsFeedNoAlerts = `<?xml version = '1.0' encoding = 'UTF-8' standalone = 'yes'?>

	<!--
	This atom/xml feed is an index to active advisories, watches and warnings
	issued by the National Weather Service.  This index file is not the complete
	Common Alerting Protocol (CAP) alert message.  To obtain the complete CAP
	alert, please follow the links for each entry in this index.  Also note the
	CAP message uses a style sheet to convey the information in a human readable
	format.  Please view the source of the CAP message to see the complete data
	set.  Not all information in the CAP message is contained in this index of
	active alerts.
	-->

	<feed
	xmlns = 'http://www.w3.org/2005/Atom'
	xmlns:cap = 'urn:oasis:names:tc:emergency:cap:1.1'
	xmlns:ha = 'http://www.alerting.net/namespace/index_1.0'
	>

	<!-- TZN = <> -->
	<!-- TZO = <> -->
	<!-- http-date = Mon, 03 Mar 2017 15:39:19 GMT -->
	<id>https://alerts.weather.gov/cap/wwaatmget.php?x=ORZ006&amp;y=0</id>
	<generator>NWS CAP Server</generator>
	<updated>2017-03-06T15:39:19+00:00</updated>
	<author>
	<name>w-nws.webmaster@noaa.gov</name>
	</author>
	<title>Current Watches, Warnings and Advisories for Greater Portland Metro Area (ORZ006) Oregon Issued by the National Weather Service</title>
	<link href='https://alerts.weather.gov/cap/wwaatmget.php?x=ORZ006&amp;y=0'/>

	    <entry>
	    <id>https://alerts.weather.gov/cap/wwaatmget.php?x=ORZ006&amp;y=0</id>
	    <updated>2017-03-06T15:39:19+00:00</updated>
	    <author>
	    <name>w-nws.webmaster@noaa.gov</name>
	    </author>
	    <title>There are no active watches, warnings or advisories</title>
	    <link href='https://alerts.weather.gov/cap/wwaatmget.php?x=ORZ006&amp;y=0'/>
	    </entry>
	    </feed>`

	alertsFeed = `<?xml version = '1.0' encoding = 'UTF-8' standalone = 'yes'?>

<!--
This atom/xml feed is an index to active advisories, watches and warnings 
issued by the National Weather Service.  This index file is not the complete 
Common Alerting Protocol (CAP) alert message.  To obtain the complete CAP 
alert, please follow the links for each entry in this index.  Also note the 
CAP message uses a style sheet to convey the information in a human readable 
format.  Please view the source of the CAP message to see the complete data 
set.  Not all information in the CAP message is contained in this index of 
active alerts.
-->

<feed
xmlns = 'http://www.w3.org/2005/Atom'
xmlns:cap = 'urn:oasis:names:tc:emergency:cap:1.1'
xmlns:ha = 'http://www.alerting.net/namespace/index_1.0'
>

<!-- TZN = <> -->
<!-- TZO = <> -->
<!-- http-date = Fri, 03 Mar 2017 21:34:46 GMT -->
<id>https://alerts.weather.gov/cap/wwaatmget.php?x=ORZ006&amp;y=0</id>
<generator>NWS CAP Server</generator>
<updated>2017-03-03T21:34:46+00:00</updated>
<author>
<name>w-nws.webmaster@noaa.gov</name>
</author>
<title>Current Watches, Warnings and Advisories for Greater Portland Metro Area (ORZ006) Oregon Issued by the National Weather Service</title>
<link href='https://alerts.weather.gov/cap/wwaatmget.php?x=ORZ006&amp;y=0'/>

<entry>
<id>https://alerts.weather.gov/cap/wwacapget.php?x=OR125843EE9C6C.SpecialWeatherStatement.125843FBD6C0OR.PDTSPSPDT.af3c363d42ff38a82ce91ac0599b12e0</id>
<updated>2017-03-06T05:31:00-08:00</updated>
<published>2017-03-06T05:31:00-08:00</published>
<author>
<name>w-nws.webmaster@noaa.gov</name>
</author>
<title>Special Weather Statement issued March 06 at 5:31AM PST until March 06 at 4:00PM PST by NWS</title>
<link href="https://alerts.weather.gov/cap/wwacapget.php?x=OR125843EE9C6C.SpecialWeatherStatement.125843FBD6C0OR.PDTSPSPDT.af3c363d42ff38a82ce91ac0599b12e0"/>
<summary>...Numerous Snow Showers Today and Tonight... ....Another Round of Snow Expected Tuesday and Tuesday Night... A cold, unstable northwest flow pattern will continue today and tonight. A weather disturbance will also be moving through the area during this time. Numerous to widespread snow showers are expected to develop this morning or early afternoon. Some of the</summary>
<cap:event>Special Weather Statement</cap:event>
<cap:effective>2017-03-06T05:31:00-08:00</cap:effective>
<cap:expires>2017-03-06T16:00:00-08:00</cap:expires>
<cap:status>Actual</cap:status>
<cap:msgType>Alert</cap:msgType>
<cap:category>Met</cap:category>
<cap:urgency>Expected</cap:urgency>
<cap:severity>Minor</cap:severity>
<cap:certainty>Observed</cap:certainty>
<cap:areaDesc>Grande Ronde Valley; Northern Blue Mountains of Oregon; Wallowa County</cap:areaDesc>
<cap:polygon></cap:polygon>
<cap:geocode>
<valueName>FIPS6</valueName>
<value>041059 041061 041063</value>
<valueName>UGC</valueName>
<value>ORZ049 ORZ050 ORZ502</value>
</cap:geocode>
<cap:parameter>
<valueName>VTEC</valueName>
<value></value>
</cap:parameter>
</entry>

</feed>`

	alert = `
<?xml version = '1.0' encoding = 'UTF-8' standalone = 'yes'?>
<?xml-stylesheet href='https://alerts.weather.gov/cap/capatomproduct.xsl' type='text/xsl'?>

<!--
This atom/xml feed is an index to active advisories, watches and warnings 
issued by the National Weather Service.  This index file is not the complete 
Common Alerting Protocol (CAP) alert message.  To obtain the complete CAP 
alert, please follow the links for each entry in this index.  Also note the 
CAP message uses a style sheet to convey the information in a human readable 
format.  Please view the source of the CAP message to see the complete data 
set.  Not all information in the CAP message is contained in this index of 
active alerts.
-->

<alert xmlns = 'urn:oasis:names:tc:emergency:cap:1.1'>

<!-- http-date = Mon, 06 Mar 2017 01:31:00 GMT -->
<identifier>NOAA-NWS-ALERTS-OR125843EE9C6C.SpecialWeatherStatement.125843FBD6C0OR.PDTSPSPDT.af3c363d42ff38a82ce91ac0599b12e0</identifier>
<sender>w-nws.webmaster@noaa.gov</sender>
<sent>2017-03-06T05:31:00-08:00</sent>
<status>Actual</status>
<msgType>Alert</msgType>
<scope>Public</scope>
<note>Alert for Grande Ronde Valley; Northern Blue Mountains of Oregon; Wallowa County (Oregon) Issued by the National Weather Service</note>
<info>
<category>Met</category>
<event>Special Weather Statement</event>
<urgency>Expected</urgency>
<severity>Minor</severity>
<certainty>Observed</certainty>
<eventCode>
<valueName>SAME</valueName>
<value>SPS</value>
</eventCode>
<effective>2017-03-06T05:31:00-08:00</effective>
<expires>2017-03-06T16:00:00-08:00</expires>
<senderName>NWS Pendleton (Northeast Oregon and Southeast Washington)</senderName>
<headline>Special Weather Statement issued March 06 at 5:31AM PST until March 06 at 4:00PM PST by NWS Pendleton</headline>
<description>...Numerous Snow Showers Today and Tonight...
....Another Round of Snow Expected Tuesday and Tuesday Night...
A cold, unstable northwest flow pattern will continue today and
tonight. A weather disturbance will also be moving through the
area during this time. Numerous to widespread snow showers are
expected to develop this morning or early afternoon. Some of the
snow showers could be heavy at times, creating snow covered roads
and limited visibility. Snow accumulations for today and tonight
are expected to range from 1 to 3 inches in the Simcoe
Highlands, Central Oregon, Ochoco John Day Highlands, Southern
Blue Mountains, Grande Ronde Valley and Wallowa Valley. Expect
higher amounts, between 2 to 5 inches along the East Slopes of the
Washington Cascades and in the Northern Blue Mountains.
Expect a brief lull in the snow shower activity late this evening
and into the first part of tonight. Then, the next storm system
will move into the area for Tuesday, Tuesday night and Wednesday
morning. This storm system will have more moisture to work with
and therefore am expecting widespread steady snow for the area
mountains. Snow levels will be rising Tuesday afternoon and
evening...reaching 4500 to 5000 feet over Central Oregon...4000
feet over Northeast Oregon and 2000 to 2500 feet in Southern
Washington. Preliminary additional snowfall totals for Tuesday
through Wednesday morning are expected to range from less than 1
inch in Central Oregon...1 to 2 inches in the Grande Ronde Valley
and Ochoco John Day Highlands...2 to 5 inches in the Simcoe
Highlands, Southern Blue Mountains and much of Wallowa County.
Expect much higher amounts between 6 to 13 inches for the East
Slopes of the Washington Cascades and the Northern Blue
Mountains where a winter storm watch is in effect. The
precipitation should gradually taper off and end for most
locations by late Tuesday night or Wednesday morning.
Those planning on traveling over mountain passes or doing
activities in the mountains should plan and prepare for these
significant snow accumulations. This is a complex late winter
weather pattern. Additional winter weather highlights may be
needed. Please check back for the latest forecast updates.</description>
<instruction></instruction>
<parameter>
<valueName>WMOHEADER</valueName>
<value></value>
</parameter>
<parameter>
<valueName>UGC</valueName>
<value>ORZ049-050-502-503-506-509-511-WAZ030-520-521</value>
</parameter>
<parameter>
<valueName>VTEC</valueName>
<value></value>
</parameter>
<parameter>
<valueName>TIME...MOT...LOC</valueName>
<value></value>
</parameter>
<area>
<areaDesc>Grande Ronde Valley; Northern Blue Mountains of Oregon; Wallowa County</areaDesc>
<polygon></polygon>
<geocode>
<valueName>FIPS6</valueName>
<value>041059</value>
</geocode>
<geocode>
<valueName>FIPS6</valueName>
<value>041061</value>
</geocode>
<geocode>
<valueName>FIPS6</valueName>
<value>041063</value>
</geocode>
<geocode>
<valueName>UGC</valueName>
<value>ORZ049</value>
</geocode>
<geocode>
<valueName>UGC</valueName>
<value>ORZ050</value>
</geocode>
<geocode>
<valueName>UGC</valueName>
<value>ORZ502</value>
</geocode>
</area>
</info>
</alert>
`
)

func RegisterResponders() {
	// HTML forecast
	httpmock.RegisterResponder("GET",
		"http://forecast.weather.gov/MapClick.php",
		func(req *http.Request) (*http.Response, error) {
			if req.URL.Query().Get("lat") == "45.53" && req.URL.Query().Get("lon") == "-122.67" {
				return httpmock.NewStringResponse(200, htmlForecast), nil
			}
			// Annoyingly, the NWS website returns 200 if the coordinates are
			// bad. The body in this case doesn't matter so is an empty string.
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	// Observation
	httpmock.RegisterResponder("GET",
		"http://w1.weather.gov/xml/current_obs/KPDX.xml",
		httpmock.NewStringResponder(200, observation),
	)

	// Alerts Feed
	httpmock.RegisterResponder("GET",
		"https://alerts.weather.gov/cap/wwaatmget.php",
		func(req *http.Request) (*http.Response, error) {
			if req.URL.Query().Get("x") == "ORZ006NoAlerts" && req.URL.Query().Get("y") == "0" {
				return httpmock.NewStringResponse(200, alertsFeedNoAlerts), nil
			}
			if req.URL.Query().Get("x") == "ORZ006" && req.URL.Query().Get("y") == "0" {
				return httpmock.NewStringResponse(200, alertsFeed), nil
			}
			// Annoyingly, the NWS website returns 200 if parameters are bad.
			// The body in this case doesn't matter so is an empty string.
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	// Alert
	httpmock.RegisterResponder("GET",
		"https://alerts.weather.gov/cap/wwacapget.php",
		func(req *http.Request) (*http.Response, error) {
			if req.URL.Query().Get("x") == "OR125843EE9C6C.SpecialWeatherStatement.125843FBD6C0OR.PDTSPSPDT.af3c363d42ff38a82ce91ac0599b12e0" {
				return httpmock.NewStringResponse(200, alert), nil
			}
			// TODO: bad parameters return an alert document indicating that the
			// alert has expired. Handle that.
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	// https://alerts.weather.gov/cap/wwacapget.php?x=OR125843EE9C6C.SpecialWeatherStatement.125843FBD6C0OR.PDTSPSPDT.af3c363d42ff38a82ce91ac0599b12e0
}
