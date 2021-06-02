import React from 'react';
import SwipeableViews from 'react-swipeable-views';
import { makeStyles, useTheme } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Typography from '@material-ui/core/Typography';
import Box from '@material-ui/core/Box';
import ReactDOM from 'react-dom';
import Toolbar from '@material-ui/core/Toolbar';
import Button from '@material-ui/core/Button';
import LoginPostForm from './loginForm';
import RigsterPostForm from './registerForm';

interface TabPanelProps {
  children: React.ReactNode;
  index: any;
  value: any;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`full-width-tabpanel-${index}`}
      aria-labelledby={`full-width-tab-${index}`}
    >
      {value === index && (
        <Box p={3}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

const useStyles = makeStyles((theme) => ({
  tab: {
    backgroundColor: theme.palette.background.paper,
    width: 500,
  },
  header: {
    flexGrow: 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
  grow: {
    flexGrow: 1,
  },
}));

function LoginAndRegisterSection() {
  const classes = useStyles();
  const theme = useTheme();
  const [value, setValue] = React.useState(0);

  const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setValue(newValue);
  };

  const handleChangeIndex = (index: number) => {
    setValue(index);
  };

  return (
    <div className={classes.tab}>
      <AppBar position="static" color="default">
        <Tabs
          value={value}
          onChange={handleChange}
          indicatorColor="primary"
          textColor="primary"
          variant="fullWidth"
          aria-label="full width tabs example"
        >
          <Tab
            label="ログイン"
            id="full-width-tab-0"
            aria-controls="full-width-tabpanel-0"
          />
          <Tab
            label="登録"
            id="full-width-tab-1"
            aria-controls="full-width-tabpanel-1"
          />
        </Tabs>
      </AppBar>
      <SwipeableViews
        axis={theme.direction === 'rtl' ? 'x-reverse' : 'x'}
        index={value}
        onChangeIndex={handleChangeIndex}
      >
        <TabPanel value={value} index={0}>
          <LoginPostForm />
        </TabPanel>
        <TabPanel value={value} index={1}>
          <RigsterPostForm />
        </TabPanel>
      </SwipeableViews>
    </div>
  );
}

function LoginAppBar() {
  const classes = useStyles();

  return (
    <div className={classes.header}>
      <AppBar position="static">
        <Toolbar>
          <Button
            color="inherit"
            onClick={() => {
              window.location.href = '/';
            }}
          >
            <Typography variant="h6" className={classes.title}>
              Minutes Application
            </Typography>
          </Button>
        </Toolbar>
      </AppBar>
    </div>
  );
}

// webpackでバンドルしている関係で存在していないIDが指定される場合がある
// エラーをそのままにしておくと、エラー以後のレンダリングがされない
if (document.getElementById('LoginAndRegister') != null) {
  ReactDOM.render(
    <LoginAndRegisterSection />,
    document.getElementById('LoginAndRegister')
  );
}
if (document.getElementById('loginHeader') != null) {
  ReactDOM.render(<LoginAppBar />, document.getElementById('loginHeader'));
}
