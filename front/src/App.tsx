import Container from "@material-ui/core/Container";
import Box from "@material-ui/core/Box";
import {
  createStyles,
  makeStyles,
  Typography,
  Theme,
  Grid,
  Card,
  CardContent,
  CardActions,
  Link,
} from "@material-ui/core";

import { getUrlInfo, validURL } from "./services/services";
import { UrlEntry } from "./components/UrlEntry";
import { useState } from "react";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      padding: theme.spacing(3, 2),
      display: "flex",
      justifyContent: "center",
    },
  })
);

const App: React.FC = () => {
  const classes = useStyles();
  const [url, setUrl] = useState("https://example.com/");
  const [validUrl, setValidUrl] = useState(true);
  const [shortUrl, setShortUrl] = useState<string>();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit() {
    setShortUrl(undefined);
    setError("");
    setLoading(true);
    try {
      const response = await getUrlInfo(url);
      console.log(response);
      setUrl(response.data.url);
      setShortUrl(response.data.redirect_url);
      setLoading(false);
    } catch (err) {
      setError(err.toString());
      setLoading(false);
    }
  }

  const handleChange = (url: string) => {
    setError("");
    setLoading(false);
    setUrl(url);
    setShortUrl(undefined);
    setValidUrl(validURL(url));
  };

  return (
    <Grid container className={classes.root} alignItems="center">
      <Grid item xs={12} md={6} lg={4}>
        <Box>
          <UrlEntry
            url={url}
            handleChange={handleChange}
            handleSubmit={handleSubmit}
            validUrl={validUrl}
            loading={loading}
          />
        </Box>
        {shortUrl ? (
          <Box>
            <Box paddingTop={5}>
              <Card>
                <CardContent>
                  <Typography>
                    <Link href={shortUrl} color="inherit">
                      {shortUrl}
                    </Link>
                  </Typography>
                </CardContent>
              </Card>
            </Box>
          </Box>
        ) : error ? (
          <Box mt={5}>
            <Typography>{error}</Typography>
          </Box>
        ) : (
          ""
        )}
      </Grid>
    </Grid>
  );
};

export default App;
