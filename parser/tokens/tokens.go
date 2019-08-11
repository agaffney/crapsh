package tokens

const (
	TOKEN_NULL = iota
	TOKEN_WORD
	TOKEN_ASSIGNMENT_WORD
	TOKEN_NAME
	TOKEN_NEWLINE
	TOKEN_IO_NUMBER
	// Operators
	TOKEN_AND_IF    // &&
	TOKEN_OR_IF     // ||
	TOKEN_DSEMI     // ;;
	TOKEN_DLESS     // <<
	TOKEN_DGREAT    // >>
	TOKEN_LESSAND   // <&
	TOKEN_GREATAND  // >&
	TOKEN_LESSGREAT // <>
	TOKEN_DLESSDASH // <<-
	TOKEN_CLOBBER   // >|
	TOKEN_SEMI      // ;
	TOKEN_PIPE      // |
	TOKEN_LESS      // <
	TOKEN_GREAT     // >
	TOKEN_AND       // &
	TOKEN_LPAREN    // (
	TOKEN_RPAREN    // )
	// Reserved words
	TOKEN_IF     // if
	TOKEN_THEN   // then
	TOKEN_ELSE   // else
	TOKEN_ELIF   // elif
	TOKEN_FI     // fi
	TOKEN_DO     // do
	TOKEN_DONE   // done
	TOKEN_CASE   // case
	TOKEN_ESAC   // esac
	TOKEN_WHILE  // while
	TOKEN_UNTIL  // until
	TOKEN_FOR    // for
	TOKEN_LBRACE // {
	TOKEN_RBRACE // }
	TOKEN_BANG   // !
	TOKEN_IN     // in
)
