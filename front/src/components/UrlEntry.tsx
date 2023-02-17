import {
  CircularProgress,
  createStyles,
  IconButton,
  InputBase,
  makeStyles,
  Paper,
  Theme,
} from "@material-ui/core";
import SearchIcon from "@material-ui/icons/Search";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      padding: "2px 4px",
      display: "flex",
      alignItems: "center",
    },
    input: {
      marginLeft: theme.spacing(1),
      flex: 1,
    },
    iconButton: {
      padding: 10,
    },
  })
);

type Props = {
  url: string;
  handleChange: (url: string) => void;
  handleSubmit: () => void;
  validUrl: boolean;
  loading: boolean;
};

export const UrlEntry: React.FC<Props> = ({
  url,
  handleChange,
  handleSubmit,
  validUrl,
  loading,
}) => {
  const classes = useStyles();

  return (
    <Paper className={classes.root}>
      <InputBase
        className={classes.input}
        placeholder="Enter URL"
        inputProps={{ "aria-label": "enter url" }}
        value={url}
        onChange={(e) => handleChange(e.target.value)}
        onKeyUp={(e) => {
          if (validUrl && e.key == "Enter") handleSubmit();
        }}
      />
      <IconButton
        className={classes.iconButton}
        aria-label="search"
        onClick={handleSubmit}
        disabled={!validUrl}
      >
        {loading ? <CircularProgress size={24} /> : <SearchIcon />}
      </IconButton>
    </Paper>
  );
};
