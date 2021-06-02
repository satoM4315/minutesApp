import React, { useState } from 'react';
import PropTypes from 'prop-types';
import DeleteIcon from '@material-ui/icons/Delete';
import IconButton from '@material-ui/core/IconButton';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';

// メッセージ削除のAPIへのURL
// eslint-disable-next-line @typescript-eslint/naming-convention
const API_URL_DELETE_MESSAGE = '/delete_message';

export default function DeleteMessageDialog(props) {
  const { onSubmitSuccessful, targetMessage, id, isHidden } = props;

  const [message, setMessage] = React.useState<string>(targetMessage);

  // Dialogを開くかどうか
  const [open, setOpen] = React.useState<boolean>(false);

  // サーバがへメッセージ追加のリクエストを処理中ならtrue、でないならfalseの状態
  const [working, setWorking] = useState<boolean>(false);

  const handleDelete = async (event: React.FormEvent) => {
    // FIXME もしかしたら、非同期なため、これが効く前にボタンをクリックできるかもしれない
    setWorking(true);
    try {
      // ページが更新されないようにする
      event.preventDefault();

      // Reactのハンドラはasyncにできる
      const res = await fetch(API_URL_DELETE_MESSAGE, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        // 相応しくないかも
        // same-originを使うべき？
        credentials: 'include',
        body: JSON.stringify({ id }),
      });
      const obj = await res.json();
      if ('error' in obj) {
        // サーバーからエラーが返却された
        throw new Error(
          `An error occurred on querying ${API_URL_DELETE_MESSAGE}, the response included error message: ${obj.error}`
        );
      }
      if (!('success' in obj)) {
        // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
        throw new Error(
          `An response from ${API_URL_DELETE_MESSAGE} unexpectedly did not have 'success' member`
        );
      }
      if (obj.success !== true) {
        throw new Error(
          `An response from ${API_URL_DELETE_MESSAGE} returned non true value as 'success' member`
        );
      }
      // 要求は成功
      // リスナ関数を呼ぶ
      onSubmitSuccessful();
    } finally {
      setWorking(false);
      setMessage('');
      setOpen(false);
    }
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  if (isHidden) {
    return null;
  }
  return (
    <span>
      <IconButton
        aria-label="delete"
        onClick={handleClickOpen}
        data-num="100"
        disabled={working}
      >
        <DeleteIcon fontSize="small" />
      </IconButton>
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          この文章を削除しますか？
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            {message}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="inherit" size="small" autoFocus>
            Cancel
          </Button>
          <Button
            onClick={handleDelete}
            color="secondary"
            size="small"
            endIcon={<DeleteIcon fontSize="small" />}
            variant="contained"
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </span>
  );
}

DeleteMessageDialog.propTypes = {
  // メッセージの更新が正常に完了したら呼ばれる関数
  onSubmitSuccessful: PropTypes.func,
  targetMessage: PropTypes.string,
  id: PropTypes.string,
  isHidden: PropTypes.bool,
};

DeleteMessageDialog.defaultProps = {
  onSubmitSuccessful: () => {
    window.location.href = window.location.pathname;
  },
  targetMessage: '',
  id: '',
  isHidden: true,
};
