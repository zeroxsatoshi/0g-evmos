package types

// Module event types
const (
	EventTypeProposalSubmit = "proposal_submit"
	EventTypeProposalClose  = "proposal_close"
	EventTypeProposalVote   = "proposal_vote"
	EventTypeRegisterVoter  = "register_voter"

	AttributeValueCategory          = "committee"
	AttributeKeyCommitteeID         = "committee_id"
	AttributeKeyProposalID          = "proposal_id"
	AttributeKeyVotingStartHeight   = "voting_start_height"
	AttributeKeyVotingEndHeight     = "voting_end_height"
	AttributeKeyProposalCloseStatus = "status"
	AttributeKeyVoter               = "voter"
	AttributeKeyBallots             = "ballots"
	AttributeKeyPublicKey           = "public_key"
	AttributeKeyProposalOutcome     = "proposal_outcome"
	AttributeKeyProposalTally       = "proposal_tally"
)
